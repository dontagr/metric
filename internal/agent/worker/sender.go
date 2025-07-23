package worker

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/internal/agent/converter"
	"github.com/dontagr/metric/internal/agent/service"
	"github.com/dontagr/metric/internal/common/hash"
	"github.com/dontagr/metric/models"
)

type Sender struct {
	cfg     *config.Config
	stats   *service.Stats
	log     *zap.SugaredLogger
	model   ModelInterface
	url     string
	client  *http.Client
	workers int
}

func NewSender(cfg *config.Config, log *zap.SugaredLogger, stats *service.Stats, lc fx.Lifecycle) *Sender {
	s := &Sender{
		cfg:    cfg,
		stats:  stats,
		log:    log,
		client: &http.Client{},
	}

	if s.cfg.RateLimit == 0 {
		s.workers = 1
		s.model = &batchModel{}
		s.url = fmt.Sprintf("http://%s/updates/", s.cfg.HTTPBindAddress)
	} else {
		s.workers = s.cfg.RateLimit
		s.model = &singleModel{}
		s.url = fmt.Sprintf("http://%s/update/", s.cfg.HTTPBindAddress)
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go s.Handle()

			return nil
		},
	})

	return s
}

type (
	ModelInterface interface {
		GetJobs(s *Sender, jobs chan any)
	}
)

func (s *Sender) worker(jobs chan any) {
	for row := range jobs {
		body, err := s.getBody(row)
		if err != nil {
			s.log.Errorf("get body: %v", err)
			break
		}

		compressedBody, err := s.compress(body)
		if err != nil {
			s.log.Errorf("compress: %v", err)
			continue
		}

		req, err := http.NewRequest("POST", s.url, compressedBody)
		if err != nil {
			s.log.Errorf("creating request: %v", err)
			continue
		}

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		if s.cfg.Security.Key != "" {
			outHash := make(chan string)
			s.GetHash(row, outHash)
			for hashRow := range outHash {
				req.Header.Add("HashSHA256", hashRow)
			}
		}

		resp, err := s.client.Do(req)
		if err != nil {
			s.log.Errorf("sending data: %v", err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			s.log.Errorf("closing response body: %v", err)
			continue
		}
	}
}

func (s *Sender) Handle() {
	jobs := make(chan any, s.workers)
	for w := 1; w <= s.workers; w++ {
		go s.worker(jobs)
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
		model.Hash = hash.ComputeHash(s.cfg.Security.Key, model)
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
