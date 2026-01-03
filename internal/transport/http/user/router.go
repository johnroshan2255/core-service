package user

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/middleware"
	"github.com/johnroshan2255/core-service/internal/user/service"
	"github.com/johnroshan2255/core-service/internal/config"
)

// SetupRoutes adds user routes to the provided router
func SetupRoutes(router *gin.Engine, userService *service.Service) {

	userHandler := NewHandler(userService)

	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)

			users.GET("/company", userHandler.GetCompanyDetails)
			users.PUT("/company", userHandler.UpdateCompanyDetails)

			users.GET("/payment", userHandler.GetPaymentDetails)
			users.PUT("/payment", userHandler.UpdatePaymentDetails)

			users.GET("/payments/history", userHandler.GetPaymentHistory)
			users.POST("/payments/history", userHandler.CreatePaymentHistory)
		}
	}
}

func SetupRouter(userService *service.Service) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{})

	SetupRoutes(router, userService)

	return router
}

func StartHTTPServer(cfg *config.Config, userService *service.Service) {
	router := SetupRouter(userService)

	port := cfg.Port
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}
	log.Printf("HTTP server (user) running on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
