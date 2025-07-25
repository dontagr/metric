package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/metric/counter"
	"github.com/dontagr/metric/internal/server/metric/factory"
	"github.com/dontagr/metric/internal/server/metric/gauge"
	"github.com/dontagr/metric/internal/server/service"
	"github.com/dontagr/metric/internal/server/service/backup"
	"github.com/dontagr/metric/internal/server/service/event"
	"github.com/dontagr/metric/internal/server/service/interfaces"
	"github.com/dontagr/metric/internal/server/service/recovery"
)

var Service = fx.Options(
	fx.Provide(
		factory.NewMetricFactory,
		event.NewEvent,
		recovery.NewRecovery,
		fx.Annotate(
			service.NewService,
			fx.As(new(interfaces.Service)),
		),
		backup.NewBackupService,
	),
	fx.Invoke(
		counter.RegisterMetric,
		gauge.RegisterMetric,
		func(*recovery.Recovery) {},
	),
)
