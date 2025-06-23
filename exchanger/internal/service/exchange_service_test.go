package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/korolev-n/gExchange/exchanger/internal/domain"
	"github.com/korolev-n/gExchange/exchanger/internal/service"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	rates []domain.ExchangeRate
	err   error
}

func (m *mockRepo) GetLatestRates(ctx context.Context) ([]domain.ExchangeRate, error) {
	return m.rates, m.err
}

func TestExchangeService_GetRates_Success(t *testing.T) {
	mock := &mockRepo{
		rates: []domain.ExchangeRate{
			{Currency: "USD", RateToRub: 90},
			{Currency: "EUR", RateToRub: 97},
			{Currency: "RUB", RateToRub: 1},
		},
	}
	svc := service.NewExchangeService(mock)

	result, err := svc.GetRates(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, float64(90), result["USD"])
	assert.Equal(t, float64(97), result["EUR"])
	assert.Equal(t, float64(1), result["RUB"])
}

func TestExchangeService_GetRates_Error(t *testing.T) {
	mock := &mockRepo{err: errors.New("db error")}
	svc := service.NewExchangeService(mock)

	_, err := svc.GetRates(context.Background())

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
