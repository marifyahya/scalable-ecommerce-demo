package main

import (
	"log"

	"github.com/marifyahya/scalable-ecommerce-demo/api-gateway/internal/config"
	"github.com/marifyahya/scalable-ecommerce-demo/api-gateway/internal/server"
)

func main() {
	cfg := config.Load()
	app := server.New(cfg)

	log.Printf("Starting API Gateway on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting API Gateway: %v", err)
	}
}
