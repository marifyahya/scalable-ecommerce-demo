package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/marifyahya/scalable-ecommerce-demo/api-gateway/internal/config"
)

// RateLimiter returns a Fiber middleware that limits requests per IP using Redis.
// Requests carrying the bypass header are exempt from rate limiting.
func RateLimiter(cfg *config.Config) fiber.Handler {
	store := redis.New(redis.Config{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Database: 0,
		Reset:    false,
	})

	return limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Get("X-Stress-Test-Bypass") == cfg.StressTestBypassSecret
		},
		Max:        100,
		Expiration: 1 * time.Minute,
		Storage:    store,
		KeyGenerator: func(c *fiber.Ctx) string {
			rawKey := c.IP() + "-" + c.Get("User-Agent")
			hash := sha256.Sum256([]byte(rawKey))
			return hex.EncodeToString(hash[:])
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Too many requests, please try again later.",
			})
		},
	})
}
