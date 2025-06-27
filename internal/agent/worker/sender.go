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
	"github.com/dontagr/metric/internal/agent/helper"
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
	url := fmt.Sprintf("http://%s/update/", s.cfg.HTTPBindAddress)
	for {
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		fmt.Printf("sender run with %v\n", s.stats.PollCount)
		for index, mType := range EnableStats {
			body, err := s.getBody(mType, index)
			if err != nil {
				fmt.Printf("Error get body for index %s: %v\n", index, err)
				continue
			}

			compressedBody, err := s.compress(body)
			if err != nil {
				fmt.Printf("Error with compress for index %s: %v\n", index, err)
				continue
			}

			req, err := http.NewRequest("POST", url, compressedBody)
			if err != nil {
				fmt.Printf("Error creating request for %s: %v\n", mType, err)
				continue
			}

			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error sending data for %s: %v\n", mType, err)
				continue
			}
			err = resp.Body.Close()
			if err != nil {
				fmt.Printf("Error closing response body for %s: %v\n", mType, err)
				continue
			}
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

func (s *Sender) getBody(mType string, index string) (*bytes.Buffer, error) {
	val := reflect.ValueOf(*s.stats).FieldByName(index)

	model, err := s.getModel(mType, index, val)
	if err != nil {
		return nil, fmt.Errorf("error creating model for %s: %v", mType, err)
	}

	modelJSON, err := json.Marshal(model)
	if err != nil {
		return nil, fmt.Errorf("error marshaling model for %s: %v", mType, err)
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
	value, err := helper.ConvertReflectValueToFloat64(val)
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
	value, err := helper.ConvertReflectValueToInt64(val)
	if err != nil {
		return nil, err
	}

	return &models.Metrics{
		ID:    index,
		MType: mType,
		Delta: &value,
	}, nil
}
