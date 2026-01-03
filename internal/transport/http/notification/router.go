package notification

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/notification"
	"github.com/johnroshan2255/core-service/internal/config"
)

// SetupRoutes adds notification routes to the provided router
func SetupRoutes(router *gin.Engine, notificationService *notification.NotificationService) {
	notificationHandler := NewHandler(notificationService)

	api := router.Group("/api/v1")
	{
		notifications := api.Group("/notifications")
		{
			notifications.POST("/user-created", notificationHandler.HandleUserCreated)
		}
	}
}

// SetupRouter creates and configures the HTTP router with notification routes
func SetupRouter(notificationService *notification.NotificationService) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{})

	SetupRoutes(router, notificationService)

	return router
}

// StartHTTPServer starts the HTTP server for notification service
func StartHTTPServer(cfg *config.Config, service *notification.NotificationService) {
	router := SetupRouter(service)

	port := cfg.Port
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}
	log.Printf("HTTP server (notification) running on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

