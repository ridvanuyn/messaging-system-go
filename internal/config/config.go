package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port             string
	DbURL            string
	RedisURL         string
	WebhookURL       string
	AuthKey          string
	MaxContentLength int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	maxContentLength, _ := strconv.Atoi(getEnvWithDefault("MAX_CONTENT_LENGTH", "160"))

	return &Config{
		Port:             getEnvWithDefault("PORT", "8080"),
		DbURL:            getEnvWithDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/messaging?sslmode=disable"),
		RedisURL:         getEnvWithDefault("REDIS_URL", "redis://localhost:6379/0"),
		WebhookURL:       getEnvWithDefault("WEBHOOK_URL", "https://webhook.site/4c65e618-2eb0-4f70-8787-3bc8681395eb"),
		AuthKey:          getEnvWithDefault("AUTH_KEY", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo"),
		MaxContentLength: maxContentLength,
	}, nil
}

// getEnvWithDefault gets environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
