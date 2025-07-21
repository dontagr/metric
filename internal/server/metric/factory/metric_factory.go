package factory

import (
	"fmt"

	"github.com/dontagr/metric/internal/server/service/intersaces"
)

type MetricFactory struct {
	Collection map[string]intersaces.Metric
}

func NewMetricFactory() *MetricFactory {
	return &MetricFactory{
		Collection: make(map[string]intersaces.Metric),
	}
}

func (f *MetricFactory) GetMetric(name string) (intersaces.Metric, error) {
	if val, ok := f.Collection[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("metric with name %s not found", name)
}

func (f *MetricFactory) SetMetric(p intersaces.Metric) {
	f.Collection[p.GetName()] = p
}
