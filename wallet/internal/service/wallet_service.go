package service

import (
	"context"
	"errors"

	"github.com/korolev-n/gExchange/exchanger/internal/repository"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetBalance(ctx context.Context, userID int64) (map[string]float64, error) {
	wallets, err := s.repo.GetBalances(ctx, userID)
	if err != nil {
		return nil, err
	}

	balances := map[string]float64{"USD": 0, "EUR": 0, "RUB": 0}
	for _, w := range wallets {
		balances[w.Currency] = w.Balance
	}
	return balances, nil
}

func (s *WalletService) Deposit(ctx context.Context, userID int64, currency string, amount float64) (map[string]float64, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	err := s.repo.UpdateBalance(ctx, userID, currency, amount)
	if err != nil {
		return nil, err
	}
	return s.GetBalance(ctx, userID)
}

func (s *WalletService) Withdraw(ctx context.Context, userID int64, currency string, amount float64) (map[string]float64, error) {
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}
	balances, err := s.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	if balances[currency] < amount {
		return nil, errors.New("insufficient funds")
	}
	err = s.repo.UpdateBalance(ctx, userID, currency, -amount)
	if err != nil {
		return nil, err
	}
	return s.GetBalance(ctx, userID)
}
