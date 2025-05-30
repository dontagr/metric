package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
)

var Config = fx.Options(
	fx.Provide(config.NewConfig),
)
