package store

import (
	"context"
	"fmt"

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

func (m *MemStorage) LoadMetric(id string, mType string) *models.Metrics {
	metrics, ok := m.collection[fmt.Sprintf("%s_%s", mType, id)]
	if !ok {
		return &models.Metrics{}
	}

	return metrics
}

func (m *MemStorage) SaveMetric(metrics *models.Metrics) {
	m.collection[fmt.Sprintf("%s_%s", metrics.MType, metrics.ID)] = metrics
}

func (m *MemStorage) ListMetric() map[string]*models.Metrics {
	return m.collection
}

func (m *MemStorage) RestoreMetricCollection(collection map[string]*models.Metrics) {
	m.collection = collection

	fmt.Printf("\u001B[032mДанные хранилища востановлены, всего метрик: %d\u001B[0m\n", len(collection))
}

func (m *MemStorage) Ping(_ context.Context) error {
	return nil
}
