package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/file"
)

var Filer = fx.Options(
	fx.Provide(file.NewFiler),
	fx.Invoke(func(*file.Filer) {}),
)
