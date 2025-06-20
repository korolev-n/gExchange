package repository

import (
	"context"
	"database/sql"

	"github.com/korolev-n/gExchange/exchanger/internal/domain"
)

type ExchangeRepository interface {
	GetLatestRates(ctx context.Context) ([]domain.ExchangeRate, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) ExchangeRepository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) GetLatestRates(ctx context.Context) ([]domain.ExchangeRate, error) {
	query := `
        SELECT set_date, currency, rate_to_rub 
        FROM exchange_rates 
        WHERE set_date = (SELECT MAX(set_date) FROM exchange_rates)
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []domain.ExchangeRate
	for rows.Next() {
		var rate domain.ExchangeRate
		if err := rows.Scan(&rate.SetDate, &rate.Currency, &rate.RateToRub); err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}
