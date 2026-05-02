package config

import "os"

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	AnthropicAPIKey string
	Port           string
}

func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://ninjo:ninjo_dev_password@db:5432/ninjo?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		AnthropicAPIKey: getEnv("ANTHROPIC_API_KEY", ""),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
