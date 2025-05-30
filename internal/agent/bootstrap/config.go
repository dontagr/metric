package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/config"
)

var Config = fx.Options(
	fx.Provide(config.NewConfig),
)
