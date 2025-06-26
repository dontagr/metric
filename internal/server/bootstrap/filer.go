package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/store"
)

var Filer = fx.Options(
	fx.Provide(store.NewFiler),
	fx.Invoke(func(*store.Filer) {}),
)
