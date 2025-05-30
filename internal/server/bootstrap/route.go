package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/service"
)

var Route = fx.Options(
	fx.Provide(
		service.NewUpdateHandler,
		service.NewServeMux,
	),
)
