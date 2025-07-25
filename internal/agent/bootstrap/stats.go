package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/service"
)

var Stats = fx.Options(
	fx.Provide(
		service.NewStats,
	),
)
