package worker

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	hash2 "hash"
	"io"
	"os"
	"reflect"
	"sync"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/internal/agent/converter"
	"github.com/dontagr/metric/internal/agent/service"
	"github.com/dontagr/metric/internal/agent/service/transport"
	"github.com/dontagr/metric/internal/common/hash"
	"github.com/dontagr/metric/models"
)

type Sender struct {
	model     ModelInterface
	transport transport.Transport
	cfg       *config.Config
	stats     *service.Stats
	log       *zap.SugaredLogger
	workers   int
	publicKey *rsa.PublicKey
	workerWG  sync.WaitGroup
	exitChan  chan bool
}

func NewSender(cfg *config.Config, log *zap.SugaredLogger, stats *service.Stats, lc fx.Lifecycle, transport *transport.HTTPManager) (*Sender, error) {
	s := &Sender{
		cfg:       cfg,
		stats:     stats,
		log:       log,
		transport: transport,
		exitChan:  make(chan bool, 1),
	}

	if s.cfg.CryptoKey != "" {
		publicKeyData, err := os.ReadFile(s.cfg.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("failed to read PEM file: %v", err)
		}

		publicKeyAny, err := parseRSAPublicKeyFromPEM(publicKeyData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %v", err)
		}

		var ok bool
		if s.publicKey, ok = publicKeyAny.(*rsa.PublicKey); !ok {
			return nil, fmt.Errorf("not an RSA public key")
		}
	}

	if s.cfg.RateLimit == 0 {
		log.Infow("Agent sender run with 1 worker and batchModel")
		s.workers = 1
		s.model = &batchModel{}
	} else {
		log.Infof("Agent sender run with %d worker and singleModel", s.cfg.RateLimit)
		s.workers = s.cfg.RateLimit
		s.model = &singleModel{}
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go s.Handle()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Infof("Получен сигнал для завершения работы. Жду оканчания отправки.")
			s.exitChan <- true
			s.workerWG.Wait()
			s.log.Infof("Завершение...")

			return nil
		},
	})

	return s, nil
}

type (
	ModelInterface interface {
		GetJobs(s *Sender, jobs chan any)
	}
)

func (s *Sender) worker(w int, jobs chan any) {
	s.log.Infof("worker %d runing", w)
	for row := range jobs {
		select {
		case <-s.exitChan:
			return
		default:
		}

		s.workerWG.Add(1)
		body, err := s.getBody(row)
		if err != nil {
			s.log.Errorf("worker %d get body: %v", w, err)
			s.workerWG.Done()
			break
		}

		compressedBody, err := s.compress(body)
		if err != nil {
			s.log.Errorf("worker %d compress: %v", w, err)
			s.workerWG.Done()
			continue
		}

		cryptoBody, err := s.crypto(compressedBody)
		if err != nil {
			s.log.Errorf("worker %d crypto: %v", w, err)
			s.workerWG.Done()
			continue
		}

		HashSHA256 := make([]string, 0, 1)
		if s.cfg.Security.Key != "" {
			outHash := make(chan string)
			s.GetHash(row, outHash)
			for hashRow := range outHash {
				HashSHA256 = append(HashSHA256, hashRow)
			}
		}

		err = s.transport.NewRequest(cryptoBody, HashSHA256, w)
		if err != nil {
			s.log.Errorf("worker %d: %v", w, err)
		}
		s.workerWG.Done()
	}
}

func (s *Sender) Handle() {
	jobs := make(chan any, s.workers)
	for w := 1; w <= s.workers; w++ {
		go s.worker(w, jobs)
	}

	for {
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)

		s.stats.UpdateWg.Wait()
		s.stats.SendWg.Add(1)
		s.model.GetJobs(s, jobs)
		s.stats.SendWg.Done()

		s.log.Infof("sender run with PollCount: %v", s.stats.PollCount)
	}
}

func (s *Sender) GetHash(row any, outHash chan<- string) {
	defer close(outHash)

	switch v := row.(type) {
	case []any:
		for _, i := range v {
			if g, ok := i.(*models.Metrics); ok {
				outHash <- g.Hash
			}
		}
		return
	case *models.Metrics:
		outHash <- v.Hash
		return
	}
}

func (s *Sender) crypto(body *bytes.Buffer) (*bytes.Buffer, error) {
	if s.publicKey == nil {
		return body, nil
	}

	encryptedBytes, err := encryptOAEP(sha256.New(), rand.Reader, s.publicKey, body.Bytes(), nil)
	if err != nil {
		return nil, fmt.Errorf("rsa.EncryptOAEP: %v", err)
	}

	return bytes.NewBuffer([]byte(base64.StdEncoding.EncodeToString(encryptedBytes))), nil
}

func encryptOAEP(hash hash2.Hash, random io.Reader, public *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func parseRSAPublicKeyFromPEM(pemBytes []byte) (any, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	if block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func (s *Sender) compress(body *bytes.Buffer) (*bytes.Buffer, error) {
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)
	defer func(gzipWriter *gzip.Writer) {
		err := gzipWriter.Close()
		if err != nil {
			s.log.Errorf("gzipWriter.Close: %v", err)
		}
	}(gzipWriter)

	_, err := gzipWriter.Write(body.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error compressing data: %w", err)
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing Gzip writer: %w", err)
	}

	return &compressedBody, nil
}

func (s *Sender) getMetric(mType string, index string) (*models.Metrics, error) {
	val := reflect.ValueOf(*s.stats).FieldByName(index)

	model, err := s.getModel(mType, index, val)
	if err != nil {
		return nil, fmt.Errorf("error creating model for %s: %w", mType, err)
	}

	if s.cfg.Security.Key != "" {
		hashManager := hash.NewHashManager()
		hashManager.SetKey(s.cfg.Security.Key)
		hashManager.SetMetrics(model)

		model.Hash = hashManager.GetHash()
	}

	return model, nil
}

func (s *Sender) getBody(body any) (*bytes.Buffer, error) {
	modelJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling body: %w", err)
	}

	return bytes.NewBuffer(modelJSON), nil
}

func (s *Sender) getModel(mType string, index string, val reflect.Value) (*models.Metrics, error) {
	if mType == models.Gauge {
		return s.getGaugeModel(mType, index, val)
	}

	return s.getCounterModel(mType, index, val)
}

func (s *Sender) getGaugeModel(mType string, index string, val reflect.Value) (*models.Metrics, error) {
	value, err := converter.ReflectValueToFloat64(val)
	if err != nil {
		return nil, err
	}

	return &models.Metrics{
		ID:    index,
		MType: mType,
		Value: &value,
	}, nil
}

func (s *Sender) getCounterModel(mType string, index string, val reflect.Value) (*models.Metrics, error) {
	value, err := converter.ReflectValueToInt64(val)
	if err != nil {
		return nil, err
	}

	return &models.Metrics{
		ID:    index,
		MType: mType,
		Delta: &value,
	}, nil
}
