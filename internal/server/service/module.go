package service

import (
	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/metric/factory"
	"github.com/dontagr/metric/internal/server/service/backup"
	"github.com/dontagr/metric/internal/store"
	"github.com/dontagr/metric/models"
)

func NewService(mf *factory.MetricFactory, sf *store.StoreFactory, cnf *config.Config, backup *backup.Service) (*Service, error) {
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
		MetricFactory: mf,
		Store:         storage,
		HashKey:       cnf.Security.Key,
		Backup:        backup,
	}

	return &s, nil
}
