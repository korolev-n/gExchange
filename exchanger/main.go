package main

import (
	"github.com/korolev-n/gExchange/exchanger/internal/config"
	"github.com/korolev-n/gExchange/exchanger/internal/logger"
	"github.com/korolev-n/gExchange/exchanger/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)
	srv := server.New(log)

	if err := srv.Start(cfg.ServerPort); err != nil {
		log.Error("Server failed", "error", err)
	}
}
