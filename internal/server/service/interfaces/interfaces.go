package interfaces

import (
	"context"

	serviceModels "github.com/dontagr/metric/internal/server/service/models"
	"github.com/dontagr/metric/models"
)

type (
	Store interface {
		LoadMetric(id string, mType string) (*models.Metrics, error)
		SaveMetric(metrics *models.Metrics) error
		BulkSaveMetric(metrics map[string]*models.Metrics) error
		ListMetric() (map[string]*models.Metrics, error)
		RestoreMetricCollection(ctx context.Context, collection map[string]*models.Metrics) error
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
	Service interface {
		GetMetric(requestMetric serviceModels.RequestMetric) (*models.Metrics, *models.InternalError)
		GetStringValue(metrics *models.Metrics) (string, *models.InternalError)
		GetAllMetricHTML() (string, *models.InternalError)
		GetAllMetrics() (map[string]*models.Metrics, *models.InternalError)
		UpdateMetrics(requestArrayMetric serviceModels.RequestArrayMetric) (map[string]*models.Metrics, *models.InternalError)
		UpdateMetric(requestMetric serviceModels.RequestMetric) (*models.Metrics, *models.InternalError)
		Ping(ctx context.Context) error
	}
	IStoreFactory interface {
		GetStore(name string) (Store, error)
		SetStory(s Store)
	}
)
