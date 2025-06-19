package server

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Server struct {
	log *slog.Logger
	db  *sql.DB
}

func New(log *slog.Logger, db *sql.DB) *Server {
	return &Server{
		log: log,
		db:  db,
	}
}

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	MethodKey    contextKey = "method"
	PathKey      contextKey = "path"
)

func (s *Server) Start(port string) error {

	s.log.Info("Starting server", "port", port)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.addContextMiddleware(http.DefaultServeMux),
	}

	return server.ListenAndServe()
}

func (s *Server) addContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		ctx = context.WithValue(ctx, MethodKey, r.Method)
		ctx = context.WithValue(ctx, PathKey, r.URL.Path)

		logger := s.log.With(
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		logger.Info("Incoming request", slog.String("event", "request_start"))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
