package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SeedInventoryData() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
    INSERT INTO exchange_rates (set_date, currency, rate_to_rub) VALUES
    ('2025-06-01', 'USD', 75.50),
    ('2025-06-01', 'EUR', 82.30),
    ('2025-06-01', 'RUB', 1.00),
    ('2025-06-15', 'USD', 76.00),
    ('2025-06-15', 'EUR', 83.00),
    ('2025-06-15', 'RUB', 1.00),
    ('2025-06-20', 'USD', 75.80),
    ('2025-06-20', 'EUR', 82.90),
    ('2025-06-20', 'RUB', 1.00);
    `

	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Printf("Error executing SQL statement: %v", err)
	}
}
