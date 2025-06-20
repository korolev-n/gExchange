package service

import (
	"context"

	"github.com/korolev-n/gExchange/exchanger/internal/repository"
)

type ExchangeService struct {
	repo repository.ExchangeRepository
}

func NewExchangeService(repo repository.ExchangeRepository) *ExchangeService {
	return &ExchangeService{repo: repo}
}

func (s *ExchangeService) GetRates(ctx context.Context) (map[string]float64, error) {
	records, err := s.repo.GetLatestRates(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	for _, r := range records {
		result[r.Currency] = r.RateToRub
	}
	return result, nil
}
