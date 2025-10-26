package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/agent/service/transport"
	crypro "github.com/dontagr/metric/pkg/crypto"
)

var Service = fx.Options(
	fx.Provide(
		transport.NewHTTPManager,
		transport.NewGRPCManager,
		crypro.NewCManager,
	),
)
