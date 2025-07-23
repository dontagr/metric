package recovery

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service/interfaces"
	"github.com/dontagr/metric/internal/store"
	"github.com/dontagr/metric/models"
)

type Recovery struct {
	store       interfaces.Store
	filer       *store.Filer
	autoRestore bool
	log         *zap.SugaredLogger
}

func NewRecovery(log *zap.SugaredLogger, sf *store.StoreFactory, filer *store.Filer, cfg *config.Config, lc fx.Lifecycle) (*Recovery, error) {
	var storeName string
	if cfg.DataBase.Init {
		storeName = models.StorePg
	} else {
		storeName = models.StoreMem
	}

	storage, err := sf.GetStore(storeName)
	if err != nil {
		return nil, err
	}

	r := Recovery{
		store:       storage,
		filer:       filer,
		autoRestore: cfg.Store.Restore,
		log:         log,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := r.ResetStoreData(ctx)
			if err != nil {
				log.Error(err)
			}

			return nil
		},
	})
	return &r, nil
}

func (r *Recovery) ResetStoreData(ctx context.Context) error {
	if !r.autoRestore {
		return nil
	}

	data, err := r.filer.Read()
	if err != nil {
		return fmt.Errorf("ошибка при восстановлении хранилища: %w", err)
	}

	var collection map[string]*models.Metrics
	err = json.Unmarshal(data, &collection)
	if err != nil {
		return fmt.Errorf("ошибка при десериализации данных из JSON: %w", err)
	}

	err = r.store.RestoreMetricCollection(ctx, collection)
	if err != nil {
		return fmt.Errorf("ошибка при записи в хранилище: %w", err)
	}

	r.log.Infof("данные хранилища востановлены, всего метрик: %d", len(collection))

	return nil
}
