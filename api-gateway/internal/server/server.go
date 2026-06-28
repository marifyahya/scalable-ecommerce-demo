package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marifyahya/scalable-ecommerce-demo/api-gateway/internal/config"
	"github.com/marifyahya/scalable-ecommerce-demo/api-gateway/internal/middleware"
)

func New(cfg *config.Config) *fiber.App {
	app := fiber.New()

	app.Use(middleware.RateLimiter(cfg))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "API Gateway is running",
		})
	})

	return app
}
