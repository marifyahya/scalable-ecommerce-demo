package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Basic health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API Gateway is running",
		})
	})

	// Get port from environment variable or default to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting API Gateway on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting API Gateway: %v", err)
	}
}
