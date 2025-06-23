package httptransport

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/korolev-n/gExchange/exchanger/internal/service"
)

type Handler struct {
	Service *service.ExchangeService
	Logger  *slog.Logger
}

func New(service *service.ExchangeService, logger *slog.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) GetRates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rates, err := h.Service.GetRates(ctx)
	if err != nil {
		h.Logger.Error("failed to get exchange rates", "error", err)

		w.WriteHeader(http.StatusInternalServerError)

		encodeErr := json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Failed to retrieve exchange rates",
		})

		if encodeErr != nil {
			h.Logger.Error("failed to encode error response", "error", encodeErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	resp := RatesResponse{
		Rates: map[string]float64{
			"USD": rates["USD"],
			"RUB": rates["RUB"],
			"EUR": rates["EUR"],
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.Logger.Error("failed to encode JSON response", "error", err)
	}
}
