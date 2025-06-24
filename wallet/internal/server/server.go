package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/korolev-n/gExchange/wallet/internal/config"
	httptransport "github.com/korolev-n/gExchange/wallet/internal/transport/http"
)

type Server struct {
	logger        *slog.Logger
	db            *sql.DB
	cfg           *config.Config
	http          *http.Server
	authHandler   *httptransport.AuthHandler
	walletHandler *httptransport.WalletHandler
}

func New(
	logger *slog.Logger,
	db *sql.DB,
	cfg *config.Config,
	authHandler *httptransport.AuthHandler,
	walletHandler *httptransport.WalletHandler,
) *Server {
	return &Server{
		logger:        logger,
		db:            db,
		cfg:           cfg,
		authHandler:   authHandler,
		walletHandler: walletHandler,
	}
}

func (s *Server) Start(port string) error {

	s.http = &http.Server{
		Addr:    ":" + port,
		Handler: s.routes(),
	}

	s.logger.Info("Starting server", "port", port)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.http.Shutdown(ctx)
}
