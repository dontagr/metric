package bootstrap

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/metric/counter"
	"github.com/dontagr/metric/internal/server/metric/factory"
	"github.com/dontagr/metric/internal/server/metric/gauge"
	"github.com/dontagr/metric/internal/server/service"
	"github.com/dontagr/metric/internal/server/service/event"
	"github.com/dontagr/metric/internal/server/service/interfaces"
	"github.com/dontagr/metric/internal/server/service/recovery"
	"github.com/dontagr/metric/internal/store"
	"github.com/dontagr/metric/models"
)

var Service = fx.Options(
	fx.Provide(
		factory.NewMetricFactory,
		event.NewEvent,
		recovery.NewRecovery,
		fx.Annotate(
			newService,
			fx.As(new(interfaces.Service)),
		),
	),
	fx.Invoke(
		counter.RegisterMetric,
		gauge.RegisterMetric,
		func(*recovery.Recovery) {},
	),
)

func newService(log *zap.SugaredLogger, mf *factory.MetricFactory, sf *store.StoreFactory, event *event.Event, cnf *config.Config, lc fx.Lifecycle) (*service.Service, error) {
	var storeName string
	if cnf.DataBase.Init {
		storeName = models.StorePg
	} else {
		storeName = models.StoreMem
	}

	storage, err := sf.GetStore(storeName)
	if err != nil {
		return nil, err
	}

	s := service.Service{
		MetricFactory:  mf,
		Store:          storage,
		Event:          event,
		IsDirectBackup: cnf.Store.Interval == 0,
		HashKey:        cnf.Security.Key,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if !s.IsDirectBackup {
				go s.AutoBackUp(cnf.Store.Interval, log)
				log.Infof("Метрики бэкапятся каждые %v секунд", cnf.Store.Interval)
			} else {
				log.Info("Метрики бэкапятся при получении")
			}

			return nil
		},
	})

	return &s, nil
}
