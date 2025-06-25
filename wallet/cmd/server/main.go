package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/korolev-n/gExchange/wallet/internal/cache"
	"github.com/korolev-n/gExchange/wallet/internal/client"
	"github.com/korolev-n/gExchange/wallet/internal/config"
	"github.com/korolev-n/gExchange/wallet/internal/db"
	"github.com/korolev-n/gExchange/wallet/internal/logger"
	"github.com/korolev-n/gExchange/wallet/internal/repository"
	"github.com/korolev-n/gExchange/wallet/internal/server"
	"github.com/korolev-n/gExchange/wallet/internal/service"
	httptransport "github.com/korolev-n/gExchange/wallet/internal/transport/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	userRepo := repository.NewUserRepository(database)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	authHandler := httptransport.NewAuthHandler(authService)

	walletRepo := repository.NewWalletRepository(database)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to exchanger service: %v", err)
	}
	defer grpcConn.Close()

	exchangerClient := client.NewExchangerClient(grpcConn)
	exchangeCache := cache.NewExchangeRateCache(30 * time.Second)
	walletService := service.NewWalletService(walletRepo, exchangerClient, exchangeCache)
	walletHandler := httptransport.NewWalletHandler(walletService)

	srv := server.New(logger, database, cfg, authHandler, walletHandler)

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
