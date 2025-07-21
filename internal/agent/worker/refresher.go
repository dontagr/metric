package worker

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/agent/config"
	"github.com/dontagr/metric/internal/agent/service"
)

type Refresher struct {
	cfg   *config.Config
	stats *service.Stats
	log   *zap.SugaredLogger
}

func NewRefresher(cfg *config.Config, log *zap.SugaredLogger, stats *service.Stats, lc fx.Lifecycle) *Refresher {
	r := &Refresher{
		cfg:   cfg,
		stats: stats,
		log:   log,
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
		s.log.Debug("refresher run")
		s.stats.Update()

		time.Sleep(time.Duration(s.cfg.PollInterval) * time.Second)
	}
}
