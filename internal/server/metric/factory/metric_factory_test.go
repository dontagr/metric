package factory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dontagr/metric/models"
)

type MockMetric struct {
	name string
}

func (m *MockMetric) GetName() string {
	return m.name
}

func (m *MockMetric) ConvertToMetrics(id string, value string) (*models.Metrics, error) {
	return &models.Metrics{}, nil
}
func (m *MockMetric) GetMetricsByData(id string, value any) (*models.Metrics, error) {
	return &models.Metrics{}, nil
}
func (m *MockMetric) Process(oldValue *models.Metrics, newValue *models.Metrics) error {
	return nil
}
func (m *MockMetric) ReturnValue(metrics *models.Metrics) string {
	return ""
}

func TestMetricFactory_GetMetric(t *testing.T) {
	factory := NewMetricFactory()

	metricName := "metric1"
	metric := &MockMetric{name: metricName}
	factory.SetMetric(metric)

	t.Run("Get existing metric", func(t *testing.T) {
		got, err := factory.GetMetric(metricName)
		assert.NoError(t, err)
		assert.Equal(t, metric, got)
	})

	t.Run("Metric not found", func(t *testing.T) {
		_, err := factory.GetMetric("unknown_metric")
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("metric with name %s not found", "unknown_metric"), err.Error())
	})
}

func TestMetricFactory_SetMetric(t *testing.T) {
	factory := NewMetricFactory()

	metricName := "metric1"
	metric := &MockMetric{name: metricName}

	factory.SetMetric(metric)

	// Проверим, что metric была успешно установлена
	got, err := factory.GetMetric(metricName)
	assert.NoError(t, err)
	assert.Equal(t, metric, got)
}
