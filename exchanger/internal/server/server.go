package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	httptransport "github.com/korolev-n/gExchange/exchanger/internal/transport/http"
)

type Server struct {
	logger  *slog.Logger
	db      *sql.DB
	http    *http.Server
	handler *httptransport.Handler
}

func New(log *slog.Logger, db *sql.DB, handler *httptransport.Handler) *Server {
	return &Server{
		logger:  log,
		db:      db,
		handler: handler,
	}
}

func (s *Server) Start(port string) error {

	s.http = &http.Server{
		Addr:    ":" + port,
		Handler: s.routes(),
	}

	//s.logger.Info("Starting server", "port", port)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.http.Shutdown(ctx)
}
