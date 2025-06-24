package service

import (
	"context"
	"errors"

	"github.com/korolev-n/gExchange/wallet/internal/repository"
)

type ExchangeRateFetcher interface {
	GetRates(ctx context.Context) (map[string]float64, error)
}

type WalletService struct {
	repo       repository.WalletRepository
	exchanger  ExchangeRateFetcher // интерфейс gRPC-клиента
}

func NewWalletService(repo repository.WalletRepository, exchanger ExchangeRateFetcher) *WalletService {
	return &WalletService{repo: repo, exchanger: exchanger}
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

func (s *WalletService) Exchange(ctx context.Context, userID int64, from, to string, amount float64) (map[string]float64, float64, error) {
	if amount <= 0 || from == to {
		return nil, 0, errors.New("invalid amount or currencies")
	}

	balances, err := s.GetBalance(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	if balances[from] < amount {
		return nil, 0, errors.New("insufficient funds")
	}

	rates, err := s.exchanger.GetRates(ctx)
	if err != nil {
		return nil, 0, err
	}

	fromRate, ok1 := rates[from]
	toRate, ok2 := rates[to]
	if !ok1 || !ok2 || fromRate == 0 {
		return nil, 0, errors.New("invalid exchange rates")
	}

	// Конвертация через рубли
	amountInRub := amount * fromRate
	exchanged := amountInRub / toRate

	// Списание и начисление
	if err := s.repo.UpdateBalance(ctx, userID, from, -amount); err != nil {
		return nil, 0, err
	}
	if err := s.repo.UpdateBalance(ctx, userID, to, exchanged); err != nil {
		return nil, 0, err
	}

	balances[to] += exchanged
	balances[from] -= amount

	return balances, exchanged, nil
}