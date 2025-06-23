package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/korolev-n/gExchange/exchanger/internal/config"
)

type Server struct {
	logger *slog.Logger
	db     *sql.DB
	http   *http.Server
	cfg    *config.Config
}

func New(log *slog.Logger, db *sql.DB, cfg *config.Config) *Server {
	return &Server{
		logger: log,
		db:     db,
		cfg:    cfg,
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
