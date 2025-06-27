package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service/event"
)

type Filer struct {
	filename string
	perm     uint32
}

func NewFiler(cfg *config.Config, event *event.Event, lc fx.Lifecycle) *Filer {
	w := Filer{
		filename: cfg.Store.FilePath + cfg.Store.FileName,
		perm:     cfg.Store.FilePerm,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go w.consumer(event.Metrics)

			return nil
		},
	})

	return &w
}

func (w *Filer) save(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Filer: Ошибка сериализации ListMetric в JSON: %v\n", err)
		return
	}

	err = os.WriteFile(w.filename, jsonData, os.FileMode(w.perm))
	if err != nil {
		fmt.Printf("Filer: Ошибка записи по адресу: %v\n", w.filename)
		return
	}

	fmt.Printf("Запись бэкапа успешна: %v\n", w.filename)
}

func (w *Filer) consumer(event <-chan interface{}) {
	for {
		data, ok := <-event
		if !ok {
			fmt.Println("Канал был закрыт. Завершение консьюмера.")
			return
		}

		w.save(data)
	}
}

func (w *Filer) Read() ([]byte, error) {
	return os.ReadFile(w.filename)
}
