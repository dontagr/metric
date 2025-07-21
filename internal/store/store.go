package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dontagr/metric/models"
)

type MemStorage struct {
	collection map[string]*models.Metrics
	name       string
}

func newMemStorage() *MemStorage {
	return &MemStorage{
		collection: make(map[string]*models.Metrics),
		name:       models.StoreMem,
	}
}

func RegisterStoreMem(ms *StoreFactory) {
	ms.SetStory(newMemStorage())
}

func (m *MemStorage) GetName() string {
	return m.name
}

func (m *MemStorage) LoadMetric(id string, mType string) (*models.Metrics, error) {
	metrics, ok := m.collection[fmt.Sprintf("%s_%s", mType, id)]
	if !ok {
		return &models.Metrics{}, pgx.ErrNoRows
	}

	return metrics, nil
}

func (m *MemStorage) SaveMetric(metrics *models.Metrics) error {
	m.collection[fmt.Sprintf("%s_%s", metrics.MType, metrics.ID)] = metrics

	return nil
}

func (m *MemStorage) BulkSaveMetric(metrics map[string]*models.Metrics) error {
	for _, metric := range metrics {
		m.collection[fmt.Sprintf("%s_%s", metric.MType, metric.ID)] = metric
	}

	return nil
}

func (m *MemStorage) ListMetric() (map[string]*models.Metrics, error) {
	return m.collection, nil
}

func (m *MemStorage) RestoreMetricCollection(_ context.Context, collection map[string]*models.Metrics) error {
	m.collection = collection

	return nil
}

func (m *MemStorage) Ping(_ context.Context) error {
	return nil
}
