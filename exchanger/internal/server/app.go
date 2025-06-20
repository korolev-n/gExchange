package server

import (
	"context"
	"log/slog"

	"github.com/korolev-n/gExchange/exchanger/internal/transport/grpc"
)

type AppServers struct {
	HTTPServer *Server
	GRPCServer *grpc.Server
	logger     *slog.Logger
}

func NewAppServers(httpServer *Server, grpcServer *grpc.Server, logger *slog.Logger) *AppServers {
	return &AppServers{
		HTTPServer: httpServer,
		GRPCServer: grpcServer,
		logger:     logger,
	}
}

func (a *AppServers) Start(httpPort string, grpcPort string) error {
	go func() {
		a.logger.Info("Starting HTTP server", "port", httpPort)
		if err := a.HTTPServer.Start(httpPort); err != nil {
			a.logger.Error("HTTP server failed to start", "error", err)
		}
	}()

	go func() {
		a.logger.Info("Starting gRPC server", "port", grpcPort)
		if err := a.GRPCServer.Start(grpcPort, a.logger); err != nil {
			a.logger.Error("gRPC server failed to start", "error", err)
		}
	}()

	return nil
}

func (a *AppServers) Shutdown(ctx context.Context) error {
	a.logger.Info("Starting graceful shutdown")

	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		a.logger.Error("HTTP server shutdown error", "error", err)
	}

	a.GRPCServer.Stop()
	a.logger.Info("gRPC server stopped")

	return nil
}
