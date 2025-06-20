package server

import (
	"database/sql"
	"log/slog"
	"net/http"
)

type Server struct {
	logger *slog.Logger
	db     *sql.DB
}

func New(log *slog.Logger, db *sql.DB) *Server {
	return &Server{
		logger: log,
		db:     db,
	}
}

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	MethodKey    contextKey = "method"
	PathKey      contextKey = "path"
)

func (s *Server) Start(port string) error {
	s.logger.Info("Starting server", "port", port)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.routes(),
	}

	return server.ListenAndServe()
}
