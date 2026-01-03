package http

import (
	"log"

	"github.com/johnroshan2255/core-service/internal/config"
)

func StartHTTPServer(cfg *config.Config, services *Services) {
	router := SetupRouter(cfg, services)

	port := cfg.Port
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("HTTP server (core-service) running on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

