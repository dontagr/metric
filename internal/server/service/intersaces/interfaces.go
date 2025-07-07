package intersaces

import (
	"context"

	"github.com/dontagr/metric/models"
)

type (
	Store interface {
		LoadMetric(id string, mType string) (*models.Metrics, error)
		SaveMetric(metrics *models.Metrics) error
		BulkSaveMetric(metrics map[string]*models.Metrics) error
		ListMetric() (map[string]*models.Metrics, error)
		RestoreMetricCollection(collection map[string]*models.Metrics) error
		GetName() string
		Ping(ctx context.Context) error
	}
	Metric interface {
		GetName() string
		ConvertToMetrics(id string, value string) (*models.Metrics, error)
		GetMetricsByData(id string, value any) (*models.Metrics, error)
		Process(oldValue *models.Metrics, newValue *models.Metrics) error
		ReturnValue(metrics *models.Metrics) string
	}
)
