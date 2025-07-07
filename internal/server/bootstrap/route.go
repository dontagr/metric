package bootstrap

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service"
	"github.com/dontagr/metric/internal/server/service/event"
	"github.com/dontagr/metric/internal/store"
	"github.com/dontagr/metric/models"
)

var Route = fx.Options(
	fx.Provide(newUpdateHandler),
	fx.Invoke(
		service.BindRoutes,
		func(*service.UpdateHandler) {},
	),
)

func newUpdateHandler(log *zap.SugaredLogger, mf *service.MetricFactory, sf *store.StoreFactory, event *event.Event, cnf *config.Config, lc fx.Lifecycle) (*service.UpdateHandler, error) {
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

	uh := service.UpdateHandler{
		MetricFactory:  mf,
		Store:          storage,
		Event:          event,
		IsDirectBackup: cnf.Store.Interval == 0,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if !uh.IsDirectBackup {
				go uh.AutoBackUp(cnf.Store.Interval, log)
				log.Infof("Метрики бэкапятся каждые %v секунд", cnf.Store.Interval)
			} else {
				log.Info("Метрики бэкапятся при получении")
			}

			return nil
		},
	})

	return &uh, nil
}
