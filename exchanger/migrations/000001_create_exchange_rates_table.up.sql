CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    set_date DATE NOT NULL,
    currency VARCHAR(3) NOT NULL CHECK (currency IN ('USD', 'EUR', 'RUB')),
    rate_to_rub FLOAT8 NOT NULL
);

CREATE UNIQUE INDEX idx_unique_rate ON exchange_rates(set_date, currency);