package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service/event"
)

type Filer struct {
	mx       sync.RWMutex
	filename string
	perm     uint32
	log      *zap.SugaredLogger
}

func NewFiler(log *zap.SugaredLogger, cfg *config.Config, event *event.Event, lc fx.Lifecycle) *Filer {
	w := Filer{
		filename: cfg.Store.FilePath + cfg.Store.FileName,
		perm:     cfg.Store.FilePerm,
		log:      log,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go w.consumer(event.Metrics)

			return nil
		},
	})

	return &w
}

func (w *Filer) save(data interface{}) error {
	w.mx.Lock()
	defer w.mx.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("filer: ошибка сериализации ListMetric в JSON: %w", err)
	}

	err = os.WriteFile(w.filename, jsonData, os.FileMode(w.perm))
	if err != nil {
		return fmt.Errorf("filer: ошибка записи по адресу: %v", w.filename)
	}

	w.log.Infof("запись бэкапа успешна. file: %v", w.filename)

	return nil
}

func (w *Filer) consumer(event <-chan interface{}) {
	for {
		data, ok := <-event
		if !ok {
			w.log.Info("канал был закрыт. Завершение консьюмера.")
			return
		}

		err := w.save(data)
		if err != nil {
			w.log.Error(err)
		}
	}
}

func (w *Filer) Read() ([]byte, error) {
	return os.ReadFile(w.filename)
}
