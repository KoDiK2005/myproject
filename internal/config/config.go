package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	Port           string
	JWTSecret      string
	AllowedOrigins []string
}

func Load() (Config, error) {
	// Загружаем .env если он есть (в продакшне его нет — ошибку игнорируем)
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		Port:           os.Getenv("PORT"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AllowedOrigins: parseOrigins(os.Getenv("ALLOWED_ORIGINS")),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func parseOrigins(raw string) []string {
	if raw == "" {
		return []string{"http://localhost:3000", "http://localhost:5173"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}
