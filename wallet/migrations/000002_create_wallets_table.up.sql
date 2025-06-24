CREATE TABLE IF NOT EXISTS wallets (
    user_id INTEGER NOT NULL REFERENCES users(id),
    currency VARCHAR(3) NOT NULL,
    balance NUMERIC NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, currency)
);