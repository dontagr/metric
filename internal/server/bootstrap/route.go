package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/config"
	"github.com/dontagr/metric/internal/server/service/handler"
	"github.com/dontagr/metric/internal/server/service/interfaces"
)

var Route = fx.Options(
	fx.Provide(newUpdateHandler),
	fx.Invoke(
		handler.BindRoutes,
		func(*handler.UpdateHandler) {},
	),
)

func newUpdateHandler(cnf *config.Config, service interfaces.Service) (*handler.UpdateHandler, error) {
	uh := handler.UpdateHandler{
		Service: service,
		HashKey: cnf.Security.Key,
	}

	return &uh, nil
}
