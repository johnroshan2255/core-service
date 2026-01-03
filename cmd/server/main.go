package main

import (
	"log"
	"sync"

	"github.com/johnroshan2255/core-service/internal/config"
	"github.com/johnroshan2255/core-service/internal/notification"
	httptransport "github.com/johnroshan2255/core-service/internal/transport/http/notification"
	grpctransport "github.com/johnroshan2255/core-service/internal/transport/grpc/notification"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg := config.LoadConfig()

	provider := cfg.NotificationProvider
	if provider == "" {
		provider = "email"
	}
	notificationFactory, err := notification.NewFactory(provider)
	if err != nil {
		log.Fatalf("Failed to create notification factory: %v", err)
	}
	notificationService := notificationFactory.NewService()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		httptransport.StartHTTPServer(cfg, notificationService)
	}()

	go func() {
		defer wg.Done()
		grpctransport.StartGRPCServer(cfg, notificationService)
	}()

	wg.Wait()
}
