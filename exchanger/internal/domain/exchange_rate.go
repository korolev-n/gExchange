package domain

import "time"

type ExchangeRate struct {
	SetDate   time.Time
	Currency  string
	RateToRub float64
}
