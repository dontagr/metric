package service

import (
	"github.com/dontagr/metric/models"
)

type (
	Store interface {
		LoadMetric(id string, mType string) *models.Metrics
		SaveMetric(metrics *models.Metrics)
		ListMetric() map[string]*models.Metrics
		RestoreMetricCollection(collection map[string]*models.Metrics)
	}
	Metric interface {
		GetName() string
		ConvertToMetrics(id string, value string) (*models.Metrics, error)
		GetMetricsByData(id string, value any) (*models.Metrics, error)
		Process(oldValue *models.Metrics, newValue *models.Metrics) error
		ReturnValue(metrics *models.Metrics) string
	}
)
