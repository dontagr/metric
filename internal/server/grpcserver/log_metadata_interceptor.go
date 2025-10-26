package grpcserver

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogMetaDataInterceptor struct {
	log *zap.SugaredLogger
}

func (l LogMetaDataInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok {
			l.log.Infof("Method %s took %s and has status %d with message '%s'", info.FullMethod, duration, grpcStatus.Code(), grpcStatus.Message())
			if grpcStatus.Code() == codes.NotFound || grpcStatus.Code() == codes.InvalidArgument {
				return resp, err
			}
		} else {
			l.log.Infof("Method %s took %s to error", info.FullMethod, duration)
			l.log.Error(err.Error())
		}
	} else {
		l.log.Infof("Method %s took %s to completed successfully", info.FullMethod, duration)
	}

	return resp, err
}

func (l LogMetaDataInterceptor) Stream(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok {
			if grpcStatus.Code() == codes.NotFound {
				return err
			}
			l.log.Error(grpcStatus.Message())
		} else {
			l.log.Error(err.Error())
		}
	}

	return err
}
