package middleware

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingServiceKey = errors.New("service key is missing in request metadata")
	ErrInvalidServiceKey = errors.New("invalid service key")
)

// BackendAuthInterceptor validates service keys for inter-service gRPC communication
// It checks for "service-key" in the gRPC metadata and validates it against the configured service key
type BackendAuthInterceptor struct {
	serviceKey string
}

// NewBackendAuthInterceptor creates a new backend authentication interceptor
func NewBackendAuthInterceptor(serviceKey string) *BackendAuthInterceptor {
	return &BackendAuthInterceptor{
		serviceKey: serviceKey,
	}
}

// UnaryInterceptor is a gRPC unary interceptor that validates service keys
func (i *BackendAuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Printf("BackendAuthInterceptor: No metadata found in request")
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		serviceKeys := md.Get("service-key")
		if len(serviceKeys) == 0 || serviceKeys[0] == "" {
			log.Printf("BackendAuthInterceptor: Service key missing in request")
			return nil, status.Errorf(codes.Unauthenticated, ErrMissingServiceKey.Error())
		}

		if serviceKeys[0] != i.serviceKey {
			log.Printf("BackendAuthInterceptor: Invalid service key provided")
			return nil, status.Errorf(codes.Unauthenticated, ErrInvalidServiceKey.Error())
		}

		return handler(ctx, req)
	}
}

// StreamInterceptor is a gRPC stream interceptor that validates service keys
func (i *BackendAuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			log.Printf("BackendAuthInterceptor: No metadata found in stream request")
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		serviceKeys := md.Get("service-key")
		if len(serviceKeys) == 0 || serviceKeys[0] == "" {
			log.Printf("BackendAuthInterceptor: Service key missing in stream request")
			return status.Errorf(codes.Unauthenticated, ErrMissingServiceKey.Error())
		}

		if serviceKeys[0] != i.serviceKey {
			log.Printf("BackendAuthInterceptor: Invalid service key provided in stream")
			return status.Errorf(codes.Unauthenticated, ErrInvalidServiceKey.Error())
		}

		return handler(srv, ss)
	}
}

