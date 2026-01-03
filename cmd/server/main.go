package main

import (
	"log"
	"sync"

	"github.com/johnroshan2255/core-service/internal/config"
	"github.com/johnroshan2255/core-service/internal/notification"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	notificationFactory, err := notification.NewFactory(cfg.Notification.Provider)
	if err != nil {
		log.Fatalf("Failed to create notification factory: %v", err)
	}
	notificationService := notificationFactory.NewService()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		notification.StartHTTPServer(cfg.Server.Host, cfg.Server.Port, notificationService)
	}()

	go func() {
		defer wg.Done()
		notification.StartGRPCServer(cfg, notificationService)
	}()

	wg.Wait()
}
