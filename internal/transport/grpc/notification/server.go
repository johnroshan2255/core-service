package notification

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/johnroshan2255/core-service/internal/config"
	"github.com/johnroshan2255/core-service/internal/middleware"
	"github.com/johnroshan2255/core-service/internal/notification"
	notificationv1 "github.com/johnroshan2255/core-service/proto/notification/v1"
)

// NewServer creates and configures the gRPC server with TLS
func NewServer(cfg *config.Config) (*grpc.Server, net.Listener, error) {
	authInterceptor := middleware.NewBackendAuthInterceptor(cfg.ServiceKey)

	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor()),
	}

	if cfg.GRPC.TLSEnabled {
		cert, err := tls.LoadX509KeyPair(cfg.GRPC.TLSCertFile, cfg.GRPC.TLSKeyFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load TLS certificates: %w", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
			MinVersion:   tls.VersionTLS12,
		}

		creds := credentials.NewTLS(tlsConfig)
		serverOpts = append(serverOpts, grpc.Creds(creds))
		log.Printf("gRPC server configured with TLS")
	} else {
		log.Printf("WARNING: gRPC server running without TLS (not recommended for production)")
	}

	grpcServer := grpc.NewServer(serverOpts...)

	grpcAddr := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen on %s: %w", grpcAddr, err)
	}

	return grpcServer, listener, nil
}

// SetupServer registers notification gRPC service on the server
func SetupServer(grpcServer *grpc.Server, service *notification.NotificationService) {
	handler := NewHandler(service)
	notificationv1.RegisterNotificationServiceServer(grpcServer, handler)
	log.Printf("Notification gRPC service registered")
}

