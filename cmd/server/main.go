package main

import (
	"context"
	"log"

	"github.com/johnroshan2255/core-service/internal/config"
	"github.com/johnroshan2255/core-service/internal/database"
	documentrepos "github.com/johnroshan2255/core-service/internal/document/repos"
	documentservice "github.com/johnroshan2255/core-service/internal/document/service"
	documentscheduler "github.com/johnroshan2255/core-service/internal/document/scheduler"
	"github.com/johnroshan2255/core-service/internal/middleware"
	"github.com/johnroshan2255/core-service/internal/notification"
	grpctransport "github.com/johnroshan2255/core-service/internal/transport/grpc/notification"
	httptransport "github.com/johnroshan2255/core-service/internal/transport/http"
	userrepos "github.com/johnroshan2255/core-service/internal/user/repos"
	userservice "github.com/johnroshan2255/core-service/internal/user/service"
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

	middleware.SetJWTKey(cfg.JWTKey)
	if cfg.JWTKey == "" {
		log.Printf("Warning: JWT key not set. JWT authentication will not be available.")
	}

	var userService *userservice.Service
	var documentService *documentservice.Service
	var expiryScheduler *documentscheduler.ExpiryScheduler

	if cfg.DBUrl != "" {
		db, err := database.InitDB(cfg)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer database.CloseDB(db)

		userRepo := userrepos.NewGORMRepository(db)
		userService = userservice.NewService(userRepo)

		documentRepo := documentrepos.NewGORMRepository(db)
		documentService = documentservice.NewService(documentRepo)

		daysBeforeExpiry := 30
		expiryScheduler = documentscheduler.NewExpiryScheduler(documentService, cfg, daysBeforeExpiry)
		expiryScheduler.Start(context.Background())
		defer expiryScheduler.Stop()
		defer expiryScheduler.Close()
	} else {
		log.Printf("Warning: DBUrl not set. User and Document services will not be available.")
	}

	go func() {
		grpctransport.StartGRPCServer(cfg, notificationService)
	}()

	services := &httptransport.Services{
		NotificationService: notificationService,
		UserService:         userService,
		DocumentService:     documentService,
	}

	httptransport.StartHTTPServer(cfg, services)
}
