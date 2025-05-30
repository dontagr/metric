package worker

import (
	"context"
	"fmt"
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
		fmt.Println("sender run")
		fmt.Println(s.stats.PollCount)

		time.Sleep(time.Duration(s.cfg.ReportInterval) * time.Second)
	}
}
