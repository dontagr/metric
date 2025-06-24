package service

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/file"
	"github.com/dontagr/metric/models"
)

type Recovery struct {
	store       Store
	filer       *file.Filer
	autoRestore bool
}

func NewRecovery(st Store, filer *file.Filer, cfg *config.Config, lc fx.Lifecycle) *Recovery {
	r := Recovery{
		store:       st,
		filer:       filer,
		autoRestore: cfg.Store.Restore,
	}
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			r.ResetStoreData()

			return nil
		},
	})
	return &r
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
