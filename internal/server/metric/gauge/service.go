package gauge

import (
	"fmt"
	"strconv"

	"github.com/dontagr/metric/internal/server/metric/factory"
	"github.com/dontagr/metric/models"
)

type Metric struct {
	name string
}

func RegisterMetric(mf *factory.MetricFactory) {
	mf.SetMetric(&Metric{name: models.Gauge})
}

func (m *Metric) GetName() string {
	return m.name
}

func (m *Metric) GetMetricsByData(id string, value any) (*models.Metrics, error) {
	val, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("не удается преобразовать значение в float64: %v", value)
	}

	return &models.Metrics{
		ID:    id,
		MType: models.Gauge,
		Value: &val,
	}, nil
}

func (m *Metric) ConvertToMetrics(id string, value string) (*models.Metrics, error) {
	if id == "" {
		return nil, fmt.Errorf("id %v is required", id)
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("%v is not a valid int64", value)
	}

	return &models.Metrics{
		ID:    id,
		MType: models.Gauge,
		Value: &val,
	}, nil
}

func (m *Metric) Process(_ *models.Metrics, _ *models.Metrics) error {
	return nil
}

func (m *Metric) ReturnValue(metrics *models.Metrics) string {
	if metrics.Value == nil {
		return "0"
	}

	return strconv.FormatFloat(*metrics.Value, 'f', -1, 64)
}
