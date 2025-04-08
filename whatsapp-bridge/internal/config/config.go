package config

import (
	"os"
	"strings"
)

type Config struct {
	StorePath      string
	ServerPort     string
	AllowedOrigins []string
	Debug          bool
}

func New() *Config {
	return &Config{
		StorePath:      getEnvOr("STORE_PATH", "store"),
		ServerPort:     getEnvOr("PORT", "8080"),
		AllowedOrigins: strings.Split(getEnvOr("ALLOWED_ORIGINS", "https://messageai.netlify.app,http://localhost:3000"), ","),
		Debug:          getEnvOr("DEBUG", "false") == "true",
	}
}

func getEnvOr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
