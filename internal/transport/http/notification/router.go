package notification

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/notification"
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

