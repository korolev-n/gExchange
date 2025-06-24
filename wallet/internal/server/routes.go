package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httptransport "github.com/korolev-n/gExchange/wallet/internal/transport/http"
)

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", s.handleHealthz)
	r.Post("/register", s.authHandler.Register)
	r.Post("/login", s.authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(httptransport.JWTAuth(s.cfg.JWTSecret))
		r.Get("/balance", s.walletHandler.Balance)
		r.Post("/wallet/deposit", s.walletHandler.Deposit)
		r.Post("/wallet/withdraw", s.walletHandler.Withdraw)
		r.Post("/exchange", s.walletHandler.Exchange)
	})

	return r
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		s.logger.Error("Failed to write response", "error", err)
	}
}
