package store

import (
	"fmt"

	"github.com/dontagr/metric/models"
)

type MemStorage struct {
	collection map[string]*models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		collection: make(map[string]*models.Metrics),
	}
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
