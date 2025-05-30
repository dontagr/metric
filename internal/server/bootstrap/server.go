package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/httpserver"
)

var Server = fx.Options(
	fx.Provide(httpserver.NewServer),
	fx.Invoke(func(*httpserver.HTTPServer) {}),
)
