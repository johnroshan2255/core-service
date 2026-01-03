package http

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/config"
	"github.com/johnroshan2255/core-service/internal/document/service"
	"github.com/johnroshan2255/core-service/internal/notification"
	documenthttp "github.com/johnroshan2255/core-service/internal/transport/http/document"
	notificationhttp "github.com/johnroshan2255/core-service/internal/transport/http/notification"
	userhttp "github.com/johnroshan2255/core-service/internal/transport/http/user"
	userservice "github.com/johnroshan2255/core-service/internal/user/service"
)

type Services struct {
	NotificationService *notification.NotificationService
	UserService         *userservice.Service
	DocumentService     *service.Service
}

func SetupRouter(cfg *config.Config, services *Services) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{})

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	if services.NotificationService != nil {
		notificationhttp.SetupRoutes(router, services.NotificationService)
	}

	if services.UserService != nil {
		userhttp.SetupRoutes(router, services.UserService)
	}

	if services.DocumentService != nil {
		documenthttp.SetupRoutes(router, services.DocumentService)
	}

	return router
}

