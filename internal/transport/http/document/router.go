package document

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johnroshan2255/core-service/internal/document/service"
	"github.com/johnroshan2255/core-service/internal/middleware"
)

func SetupRoutes(router *gin.Engine, documentService *service.Service) {
	uploadPath := os.Getenv("DOCUMENT_UPLOAD_PATH")
	if uploadPath == "" {
		uploadPath = "./uploads/documents"
	}

	documentHandler := NewHandler(documentService, uploadPath)

	api := router.Group("/api/v1")
	{
		documents := api.Group("/documents")
		documents.Use(middleware.AuthMiddleware())
		{
			documents.POST("", documentHandler.UploadDocument)
			documents.GET("", documentHandler.ListDocuments)
			documents.GET("/:id", documentHandler.GetDocument)
			documents.PUT("/:id", documentHandler.UpdateDocument)
			documents.DELETE("/:id", documentHandler.DeleteDocument)
		}
	}
}

func SetupRouter(documentService *service.Service) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.SetTrustedProxies([]string{})

	SetupRoutes(router, documentService)

	return router
}

func StartHTTPServer(cfg interface{}, documentService *service.Service) {
	router := SetupRouter(documentService)

	port := ":8080"
	log.Printf("HTTP server (document) running on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
