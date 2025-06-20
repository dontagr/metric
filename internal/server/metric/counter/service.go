package counter

import (
	"fmt"
	"strconv"

	"github.com/dontagr/metric/internal/server/service"
	"github.com/dontagr/metric/models"
)

type Metric struct {
	name string
}

func RegisterMetric(mf *service.MetricFactory) {
	mf.SetMetric(&Metric{name: models.Counter})
}

func (m *Metric) GetName() string {
	return m.name
}

func (m *Metric) GetMetricsByData(id string, value any) (*models.Metrics, error) {
	val, ok := value.(int64)
	if !ok {
		return nil, fmt.Errorf("не удается преобразовать значение в int64: %v", value)
	}

	return &models.Metrics{
		ID:    id,
		MType: models.Counter,
		Delta: &val,
	}, nil
}

func (m *Metric) ConvertToMetrics(id string, value string) (*models.Metrics, error) {
	if id == "" {
		return nil, fmt.Errorf("id %v is required", id)
	}

	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("value %v is not a valid int64", value)
	}

	return &models.Metrics{
		ID:    id,
		MType: models.Counter,
		Delta: &val,
	}, nil
}

func (m *Metric) Process(oldValue *models.Metrics, newValue *models.Metrics) error {
	if oldValue == nil || oldValue.Delta == nil {
		return nil
	}
	if newValue == nil || newValue.Delta == nil {
		return fmt.Errorf("invalid Metrics provided: either newValue or Delta is nil")
	}

	*newValue.Delta = *newValue.Delta + *oldValue.Delta

	return nil
}

func (m *Metric) ReturnValue(metrics *models.Metrics) string {
	if metrics.Delta == nil {
		return "0"
	}

	return strconv.FormatInt(*metrics.Delta, 10)
}
