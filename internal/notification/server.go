package notification

import (
	"fmt"
	"log"

	"github.com/johnroshan2255/core-service/internal/config"
	httphandler "github.com/johnroshan2255/core-service/internal/transport/http/notification"
	grpchandler "github.com/johnroshan2255/core-service/internal/transport/grpc/notification"
)

// StartHTTPServer starts the HTTP server for notification service
func StartHTTPServer(host, port string, service *NotificationService) {
	router := httphandler.SetupRouter(service)

	addr := fmt.Sprintf("%s:%s", host, port)
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("HTTP server (frontend) running on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

// StartGRPCServer starts the gRPC server for notification service
func StartGRPCServer(cfg *config.Config, service *NotificationService) {
	grpcServer, grpcListener, err := grpchandler.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to setup gRPC server: %v", err)
	}

	grpchandler.SetupServer(grpcServer, service)

	grpcPort := fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.Port)
	if grpcPort == "" {
		grpcPort = ":9090"
	}
	log.Printf("gRPC server (backend-to-backend) running on %s", grpcPort)
	if err := grpcServer.Serve(grpcListener); err != nil {
		log.Fatalf("failed to start gRPC server: %v", err)
	}
}

