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
	"github.com/korolev-n/gExchange/exchanger/internal/repository"
	"github.com/korolev-n/gExchange/exchanger/internal/server"
	"github.com/korolev-n/gExchange/exchanger/internal/service"
	"github.com/korolev-n/gExchange/exchanger/internal/transport/grpc"
	httptransport "github.com/korolev-n/gExchange/exchanger/internal/transport/http"
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

	repo := repository.NewPostgresRepo(database)
	svc := service.NewExchangeService(repo)

	httpHandler := httptransport.New(svc, logger)
	httpServer := server.New(logger, database, httpHandler)

	grpcHandler := grpc.New(svc, logger)
	grpcServer := grpc.NewGRPCServer(grpcHandler, logger)

	app := server.NewAppServers(httpServer, grpcServer, logger)

	if err := app.Start(cfg.ServerPort, ":50051"); err != nil {
		logger.Error("Failed to start servers", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful shutdown failed", "error", err)
	}
}
