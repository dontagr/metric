package factory

import (
	"fmt"

	"github.com/dontagr/metric/internal/server/service/interfaces"
)

type MetricFactory struct {
	Collection map[string]interfaces.Metric
}

func NewMetricFactory() *MetricFactory {
	return &MetricFactory{
		Collection: make(map[string]interfaces.Metric),
	}
}

func (f *MetricFactory) GetMetric(name string) (interfaces.Metric, error) {
	if val, ok := f.Collection[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("metric with name %s not found", name)
}

func (f *MetricFactory) SetMetric(p interfaces.Metric) {
	f.Collection[p.GetName()] = p
}
