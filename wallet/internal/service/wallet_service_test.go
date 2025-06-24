package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/korolev-n/gExchange/wallet/internal/cache"
	"github.com/korolev-n/gExchange/wallet/internal/domain"
	"github.com/stretchr/testify/require"
)

type mockWalletRepo struct {
	GetBalancesFunc   func(ctx context.Context, userID int64) ([]domain.Wallet, error)
	UpdateBalanceFunc func(ctx context.Context, userID int64, currency string, delta float64) error
}

func (m *mockWalletRepo) GetBalances(ctx context.Context, userID int64) ([]domain.Wallet, error) {
	return m.GetBalancesFunc(ctx, userID)
}

func (m *mockWalletRepo) UpdateBalance(ctx context.Context, userID int64, currency string, delta float64) error {
	return m.UpdateBalanceFunc(ctx, userID, currency, delta)
}

type mockExchanger struct {
	GetRatesFunc func(ctx context.Context) (map[string]float64, error)
}

func (m *mockExchanger) GetRates(ctx context.Context) (map[string]float64, error) {
	return m.GetRatesFunc(ctx)
}

func TestWalletService_Exchange_Success(t *testing.T) {
	repo := &mockWalletRepo{
		GetBalancesFunc: func(ctx context.Context, userID int64) ([]domain.Wallet, error) {
			return []domain.Wallet{
				{UserID: userID, Currency: "USD", Balance: 100},
				{UserID: userID, Currency: "EUR", Balance: 0},
			}, nil
		},
		UpdateBalanceFunc: func(ctx context.Context, userID int64, currency string, delta float64) error {
			return nil
		},
	}

	exchanger := &mockExchanger{
		GetRatesFunc: func(ctx context.Context) (map[string]float64, error) {
			return map[string]float64{
				"USD": 1.0,
				"EUR": 1.1765,
			}, nil
		},
	}

	c := cache.NewExchangeRateCache(5 * time.Second)
	service := NewWalletService(repo, exchanger, c)

	bal, exchanged, err := service.Exchange(context.Background(), 1, "USD", "EUR", 100)
	require.NoError(t, err)
	require.InDelta(t, 85.0, exchanged, 0.01)
	require.Equal(t, 0.0, bal["USD"])
	require.InDelta(t, 85.0, bal["EUR"], 0.01)
}

func TestWalletService_Exchange_InsufficientFunds(t *testing.T) {
	repo := &mockWalletRepo{
		GetBalancesFunc: func(ctx context.Context, userID int64) ([]domain.Wallet, error) {
			return []domain.Wallet{
				{UserID: userID, Currency: "USD", Balance: 10},
			}, nil
		},
	}
	exchanger := &mockExchanger{}
	c := cache.NewExchangeRateCache(5 * time.Second)

	service := NewWalletService(repo, exchanger, c)

	_, _, err := service.Exchange(context.Background(), 1, "USD", "EUR", 100)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insufficient funds")
}

func TestWalletService_Exchange_GetRatesError(t *testing.T) {
	repo := &mockWalletRepo{
		GetBalancesFunc: func(ctx context.Context, userID int64) ([]domain.Wallet, error) {
			return []domain.Wallet{
				{UserID: userID, Currency: "USD", Balance: 100},
			}, nil
		},
	}
	exchanger := &mockExchanger{
		GetRatesFunc: func(ctx context.Context) (map[string]float64, error) {
			return nil, errors.New("gRPC failed")
		},
	}
	c := cache.NewExchangeRateCache(5 * time.Second)
	service := NewWalletService(repo, exchanger, c)

	_, _, err := service.Exchange(context.Background(), 1, "USD", "EUR", 100)
	require.Error(t, err)
	require.Contains(t, err.Error(), "gRPC failed")
}
