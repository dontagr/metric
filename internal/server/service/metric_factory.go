package service

import (
	"fmt"
)

type MetricFactory struct {
	Collection map[string]Metric
}

func NewMetricFactory() *MetricFactory {
	return &MetricFactory{
		Collection: make(map[string]Metric),
	}
}

func (f *MetricFactory) GetMetric(name string) (Metric, error) {
	if val, ok := f.Collection[name]; ok {
		return val, nil
	}

	return nil, fmt.Errorf("metric with name %s not found", name)
}

func (f *MetricFactory) SetMetric(p Metric) {
	f.Collection[p.GetName()] = p
}
