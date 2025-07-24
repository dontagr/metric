package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/service/transport"
)

var Service = fx.Options(
	fx.Provide(
		transport.NewHTTPManager,
	),
)
