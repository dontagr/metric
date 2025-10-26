package grpcserver

import (
	"context"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	metrics "github.com/dontagr/metric/grpc/gen"
	"github.com/dontagr/metric/internal/server/config"
	grpc2 "github.com/dontagr/metric/internal/server/service/grpc"
)

type GrpcServer struct {
	server *grpc.Server
}

func NewGrpcServer(cfg *config.Config, handler *grpc2.Handler, log *zap.SugaredLogger, lc fx.Lifecycle, shutdowner fx.Shutdowner) (*GrpcServer, error) {
	if cfg.GRPCServer == "" {
		log.Infof("Disabled starting gRPC server.")
		return &GrpcServer{}, nil
	}

	logInterceptor := LogMetaDataInterceptor{log: log}
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInterceptor.Unary,
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			logInterceptor.Stream,
			recovery.StreamServerInterceptor(),
		),
	)
	reflection.Register(srv)
	metrics.RegisterMetricServiceServer(srv, handler)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Infof("Starting gRPC server. Bind: %s", cfg.GRPCServer)

			ln, err := net.Listen("tcp", cfg.GRPCServer)
			if err != nil {
				return err
			}

			go func() {
				if err := srv.Serve(ln); err != nil {
					log.Error("Failed to Serve gRPC", zap.Error(err))
					_ = shutdowner.Shutdown()
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("Stopping gRPC server")

			srv.GracefulStop()

			return nil
		},
	})

	return &GrpcServer{
		server: srv,
	}, nil
}
