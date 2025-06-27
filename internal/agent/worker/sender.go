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

	"github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/internal/agent/converter"
	"github.com/dontagr/metric/internal/agent/service"
	"github.com/dontagr/metric/models"
)

type Sender struct {
	cfg   *config.Config
	stats *service.Stats
}

func NewSender(cfg *config.Config, stats *service.Stats, lc fx.Lifecycle) *Sender {
	s := &Sender{
		cfg:   cfg,
		stats: stats,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go s.Handle()

			return nil
		},
	})

	return s
}

func (s *Sender) Handle() {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/updates/", s.cfg.HTTPBindAddress)
	for {
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		fmt.Printf("sender run with %v\n", s.stats.PollCount)

		metrics := make([]*models.Metrics, 0, len(EnableStats))
		for index, mType := range EnableStats {
			metric, err := s.getMetric(mType, index)
			if err != nil {
				fmt.Printf("Error get metrics for index %s: %v\n", index, err)
				continue
			}

			metrics = append(metrics, metric)
		}

		body, err := s.getBody(metrics)
		if err != nil {
			fmt.Printf("Error get body: %v\n", err)
			continue
		}

		compressedBody, err := s.compress(body)
		if err != nil {
			fmt.Printf("Error with compress: %v\n", err)
			continue
		}

		req, err := http.NewRequest("POST", url, compressedBody)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}

		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending data: %v\n", err)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
			continue
		}
	}
}

func (s *Sender) compress(body *bytes.Buffer) (*bytes.Buffer, error) {
	var compressedBody bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBody)
	defer func(gzipWriter *gzip.Writer) {
		err := gzipWriter.Close()
		if err != nil {
			fmt.Printf("Error with gzipWriter.Close: %v\n", err)
		}
	}(gzipWriter)

	_, err := gzipWriter.Write(body.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error compressing data: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing Gzip writer: %v", err)
	}

	return &compressedBody, nil
}

func (s *Sender) getMetric(mType string, index string) (*models.Metrics, error) {
	val := reflect.ValueOf(*s.stats).FieldByName(index)

	model, err := s.getModel(mType, index, val)
	if err != nil {
		return nil, fmt.Errorf("error creating model for %s: %v", mType, err)
	}

	return model, nil
}

func (s *Sender) getBody(body []*models.Metrics) (*bytes.Buffer, error) {
	modelJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling body: %v", err)
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
