package notification

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/notification"
	"github.com/johnroshan2255/core-service/internal/config"
)

// SetupRouter creates and configures the HTTP router with notification routes
func SetupRouter(notificationService *notification.NotificationService) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{})

	notificationHandler := NewHandler(notificationService)

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		notifications := api.Group("/notifications")
		{
			notifications.POST("/user-created", notificationHandler.HandleUserCreated)
		}
	}

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
	log.Printf("HTTP server (frontend) running on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

