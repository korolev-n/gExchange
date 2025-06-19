package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL      string
	ServerPort string
	LogLevel   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBURL:      os.Getenv("DB_URL"),
		ServerPort: os.Getenv("SERVER_PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}, nil
}
