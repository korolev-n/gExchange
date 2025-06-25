package grpc

import (
	"context"
	"log/slog"

	"github.com/korolev-n/gExchange/shared/api"
	"github.com/korolev-n/gExchange/exchanger/internal/service"
)

type Handler struct {
	api.UnimplementedExchangerServiceServer
	Service *service.ExchangeService
	Logger  *slog.Logger
}

func New(service *service.ExchangeService, logger *slog.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) GetRates(ctx context.Context, _ *api.Empty) (*api.RatesResponse, error) {
	rates, err := h.Service.GetRates(ctx)
	if err != nil {
		h.Logger.Error("failed to get rates", "error", err)
		return nil, err
	}

	return &api.RatesResponse{Rates: rates}, nil
}
