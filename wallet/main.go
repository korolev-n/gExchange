package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	logger := logger.New(cfg.LogLevel)

	database, err := db.New(cfg.DBURL, logger)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}
	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("Failed to close database", "error", err)
		}
	}()

	srv := server.New(logger, database)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(cfg.ServerPort); err != nil && err.Error() != "http: Server closed" {
			logger.Error("Server failed", "error", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful shutdown failed", "error", err)
	}
}
