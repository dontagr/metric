package worker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/internal/agent/service"
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
			val := reflect.ValueOf(*s.stats).FieldByName(index)

			resp, err := http.Post(
				fmt.Sprintf("http://%s/update/%s/%s/%v", s.cfg.HttpBindAddress, mType, index, val),
				"text/plain",
				nil,
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
