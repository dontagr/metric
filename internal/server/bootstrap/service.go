package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/metric/counter"
	"github.com/dontagr/metric/internal/server/metric/gauge"
	"github.com/dontagr/metric/internal/server/service"
	"github.com/dontagr/metric/internal/server/service/event"
)

var Service = fx.Options(
	fx.Provide(
		service.NewMetricFactory,
		event.NewEvent,
		service.NewRecovery,
	),
	fx.Invoke(
		counter.RegisterMetric,
		gauge.RegisterMetric,
		func(*service.Recovery) {},
	),
)
