package repository

import (
	"context"
	"database/sql"

	"github.com/korolev-n/gExchange/wallet/internal/domain"
)

type WalletRepository interface {
	GetBalances(ctx context.Context, userID int64) ([]domain.Wallet, error)
	UpdateBalance(ctx context.Context, userID int64, currency string, amount float64) error
}

type walletRepo struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) GetBalances(ctx context.Context, userID int64) ([]domain.Wallet, error) {
	query := `SELECT currency, balance FROM wallets WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []domain.Wallet
	for rows.Next() {
		var wallet domain.Wallet
		wallet.UserID = userID
		if err := rows.Scan(&wallet.Currency, &wallet.Balance); err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (r *walletRepo) UpdateBalance(ctx context.Context, userID int64, currency string, delta float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO wallets (user_id, currency, balance)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, currency)
		DO UPDATE SET balance = wallets.balance + EXCLUDED.balance
	`, userID, currency, delta)
	if err != nil {
		return err
	}

	return tx.Commit()
}
