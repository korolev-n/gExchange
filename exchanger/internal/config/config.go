package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL      string
	ServerPort string
	LogLevel   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(filepath.Join("exchanger", ".env")); err != nil {
		return nil, fmt.Errorf("loading env: %w", err)
	}

	return &Config{
		DBURL:      os.Getenv("DB_URL"),
		ServerPort: os.Getenv("SERVER_PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}, nil
}
