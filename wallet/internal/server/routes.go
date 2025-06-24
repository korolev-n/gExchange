package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/korolev-n/gExchange/exchanger/internal/repository"
	"github.com/korolev-n/gExchange/exchanger/internal/service"
	httptransport "github.com/korolev-n/gExchange/exchanger/internal/transport/http"
)

func (s *Server) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	userRepo := repository.NewUserRepository(s.db)
	authService := service.NewAuthService(userRepo, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	authHandler := httptransport.NewAuthHandler(authService)

	walletRepo := repository.NewWalletRepository(s.db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := httptransport.NewWalletHandler(walletService)

	r.Get("/healthz", s.handleHealthz)
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(httptransport.JWTAuth(s.cfg.JWTSecret))
		r.Get("/balance", walletHandler.Balance)
		r.Post("/wallet/deposit", walletHandler.Deposit)
		r.Post("/wallet/withdraw", walletHandler.Withdraw)
	})

	return r
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		s.logger.Error("Failed to write response", "error", err)
	}
}
