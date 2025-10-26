package bootstrap

import (
	"go.uber.org/fx"

	"github.com/dontagr/metric/internal/server/grpcserver"
	"github.com/dontagr/metric/internal/server/httpserver"
	"github.com/dontagr/metric/internal/server/service/grpc"
)

var Server = fx.Options(
	fx.Provide(
		httpserver.NewServer,
		grpcserver.NewGrpcServer,
		grpc.NewHandler,
	),
	fx.Invoke(
		func(*httpserver.HTTPServer) {},
		func(*grpcserver.GrpcServer) {},
	),
)
