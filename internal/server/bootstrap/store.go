package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/service"
	"github.com/dontagr/metric/internal/store"
)

var Store = fx.Options(
	fx.Provide(
		store.NewStoreFactory,
	),
	fx.Invoke(
		store.RegisterStoreMem,
		store.RegisterStorePG,
		func(*service.Recovery) {},
	),
)
