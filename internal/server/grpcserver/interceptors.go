package grpcserver

import (
	"context"
	"net/netip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func metadataIPCheckInterceptor(trustNetwork netip.Prefix) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.PermissionDenied, "missing metadata")
		}

		xRealIP, exists := md["x-real-ip"]
		if !exists {
			return nil, status.Errorf(codes.PermissionDenied, "X-Real-IP header is required")
		}

		ip, err := netip.ParseAddr(xRealIP[0])
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "X-Real-IP failed")
		}

		if !trustNetwork.Contains(ip) {
			return nil, status.Errorf(codes.PermissionDenied, "go out busters. you ip %s out of range", xRealIP)
		}

		return handler(ctx, req)
	}
}
