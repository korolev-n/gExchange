package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Server struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Server {
	return &Server{log: log}
}

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
		ctx = context.WithValue(ctx, "request_id", requestID)
		ctx = context.WithValue(ctx, "method", r.Method)
		ctx = context.WithValue(ctx, "path", r.URL.Path)

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
