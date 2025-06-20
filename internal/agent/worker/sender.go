package worker

import (
	"bytes"
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
	for {
		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
		fmt.Printf("sender run with %v\n", s.stats.PollCount)
		for index, mType := range EnableStats {
			body, err := s.getBody(mType, index)
			if err != nil {
				fmt.Printf("Error get body for index %s\n", index)
				continue
			}

			resp, err := http.Post(
				fmt.Sprintf("http://%s/update/", s.cfg.HTTPBindAddress),
				"application/json",
				body,
			)
			if err != nil {
				fmt.Printf("Error sending data for %s: %v\n", mType, err)
				continue
			}
			err = resp.Body.Close()
			if err != nil {
				fmt.Printf("Error sending data for %s: %v\n", mType, err)
			}
		}
	}
}
func (s *Sender) getBody(mType string, index string) (*bytes.Buffer, error) {
	val := reflect.ValueOf(*s.stats).FieldByName(index)

	model, err := s.getModel(mType, index, val)
	if err != nil {
		return nil, fmt.Errorf("Error creating model for %s: %v\n", mType, err)
	}

	modelJSON, err := json.Marshal(model)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling model for %s: %v\n", mType, err)
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
