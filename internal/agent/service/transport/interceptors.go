package transport

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ipInterceptor(ip string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md := metadata.Pairs("X-Real-IP", ip)
		ctxWithMetadata := metadata.NewOutgoingContext(ctx, md)

		return invoker(ctxWithMetadata, method, req, reply, cc, opts...)
	}
}
