package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                   string
	RedisHost              string
	RedisPort              int
	JWTSecret              string
	StressTestBypassSecret string
}

// Load reads environment variables and returns a populated Config.
func Load() *Config {
	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	return &Config{
		Port:                   getEnv("PORT", "8000"),
		RedisHost:              getEnv("REDIS_HOST", "localhost"),
		RedisPort:              redisPort,
		JWTSecret:              getEnv("JWT_SECRET", ""),
		StressTestBypassSecret: getEnv("STRESS_TEST_BYPASS_SECRET", "super-secret-bypass"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
