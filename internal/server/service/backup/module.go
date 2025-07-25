package backup

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service/event"
	"github.com/dontagr/metric/internal/store"
	"github.com/dontagr/metric/models"
)

func NewBackupService(log *zap.SugaredLogger, sf *store.StoreFactory, event *event.Event, cnf *config.Config, lc fx.Lifecycle) (*Service, error) {
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

	s := Service{
		Store:          storage,
		Event:          event,
		IsDirectBackup: cnf.Store.Interval == 0,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if !s.IsDirectBackup {
				go s.autoBackUp(cnf.Store.Interval, log)
				log.Infof("Метрики бэкапятся каждые %v секунд", cnf.Store.Interval)
			} else {
				log.Info("Метрики бэкапятся при получении")
			}

			return nil
		},
	})

	return &s, nil
}
