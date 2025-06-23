package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL         string
	ServerPort    string
	LogLevel      string
	JWTSecret     string
	JWTExpiration int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading env: %w", err)
	}

	expiration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if expiration == 0 {
		expiration = 24
	}

	return &Config{
		DBURL:         os.Getenv("DB_URL"),
		ServerPort:    os.Getenv("SERVER_PORT"),
		LogLevel:      os.Getenv("LOG_LEVEL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: expiration,
	}, nil
}
