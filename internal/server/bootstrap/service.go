package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/metric/counter"
	"github.com/dontagr/metric/internal/server/metric/gauge"
	"github.com/dontagr/metric/internal/server/service"
)

var Service = fx.Options(
	fx.Provide(service.NewMetricFactory),
	fx.Invoke(
		counter.RegisterMetric,
		gauge.RegisterMetric,
	),
)
