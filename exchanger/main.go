package main

import (
	"github.com/korolev-n/gExchange/exchanger/internal/config"
	"github.com/korolev-n/gExchange/exchanger/internal/db"
	"github.com/korolev-n/gExchange/exchanger/internal/logger"
	"github.com/korolev-n/gExchange/exchanger/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)

	database, err := db.New(cfg.DBURL, log)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Error("Failed to close database", "error", err)
		}
	}()

	srv := server.New(log, database)

	if err := srv.Start(cfg.ServerPort); err != nil {
		log.Error("Server failed", "error", err)
	}
}
