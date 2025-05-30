package worker

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/internal/agent/service"
)

type Refresher struct {
	cfg   *config.Config
	stats *service.Stats
}

func NewRefresher(cfg *config.Config, stats *service.Stats, lc fx.Lifecycle) *Refresher {
	r := &Refresher{
		cfg:   cfg,
		stats: stats,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go r.Handle()

			return nil
		},
	})

	return r
}

func (s *Refresher) Handle() {
	for {
		fmt.Println("Refresher run")
		s.stats.Update()

		time.Sleep(time.Duration(s.cfg.PollInterval) * time.Second)
	}
}
