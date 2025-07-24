package backup

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/service/event"
	"github.com/dontagr/metric/internal/server/service/interfaces"
)

type Service struct {
	IsDirectBackup bool
	Store          interfaces.Store
	Event          *event.Event
	log            *zap.SugaredLogger
}

func (s *Service) Process() {
	if s.IsDirectBackup {
		metric, err := s.Store.ListMetric()
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			s.log.Errorf("error write backup %v", err)
		}
		s.Event.Metrics <- metric
	}
}

func (s *Service) autoBackUp(interval int, log *zap.SugaredLogger) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)

		metric, err := s.Store.ListMetric()
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			log.Error(err)
		} else {
			s.Event.Metrics <- metric
		}
	}
}
