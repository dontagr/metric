package service

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service/intersaces"
	"github.com/dontagr/metric/internal/store"
	"github.com/dontagr/metric/models"
)

type Recovery struct {
	store       intersaces.Store
	filer       *store.Filer
	autoRestore bool
}

func NewRecovery(sf *store.StoreFactory, filer *store.Filer, cfg *config.Config, lc fx.Lifecycle) (*Recovery, error) {
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
	}
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			r.ResetStoreData()

			return nil
		},
	})
	return &r, nil
}

func (r *Recovery) ResetStoreData() {
	if !r.autoRestore {
		return
	}

	data, err := r.filer.Read()
	if err != nil {
		fmt.Printf("Ошибка при восстановлении хранилища: %v\n", err)
		return
	}

	var collection map[string]*models.Metrics
	err = json.Unmarshal(data, &collection)
	if err != nil {
		fmt.Printf("Ошибка при десериализации данных из JSON: %v\n", err)
		return
	}

	r.store.RestoreMetricCollection(collection)
}
