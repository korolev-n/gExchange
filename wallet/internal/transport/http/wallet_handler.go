package httptransport

import (
	"encoding/json"
	"net/http"

	"github.com/korolev-n/gExchange/wallet/internal/service"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(walletService *service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) Balance(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user_claims").(*service.Claims)
	balance, err := h.walletService.GetBalance(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, `{"error":"failed to retrieve balance"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"balance": balance})
}

type operationRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user_claims").(*service.Claims)
	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	balance, err := h.walletService.Deposit(r.Context(), claims.UserID, req.Currency, req.Amount)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"message":     "Account topped up successfully",
		"new_balance": balance,
	})
}

func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user_claims").(*service.Claims)
	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	balance, err := h.walletService.Withdraw(r.Context(), claims.UserID, req.Currency, req.Amount)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"message":     "Withdrawal successful",
		"new_balance": balance,
	})
}

type exchangeRequest struct {
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Amount       float64 `json:"amount"`
}

func (h *WalletHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user_claims").(*service.Claims)

	var req exchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
		return
	}

	balances, exchanged, err := h.walletService.Exchange(r.Context(), claims.UserID, req.FromCurrency, req.ToCurrency, req.Amount)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"message":          "Exchange successful",
		"exchanged_amount": exchanged,
		"new_balance":      balances,
	})
}
