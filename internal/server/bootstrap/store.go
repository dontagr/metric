package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/service"
	"github.com/dontagr/metric/internal/server/store"
)

var Store = fx.Options(
	fx.Provide(
		fx.Annotate(
			store.NewMemStorage,
			fx.As(new(service.Store)),
		),
	),
)
