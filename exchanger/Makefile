.PHONY: run test

run:
	go run main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path=./migrations -database=$DB_URL up

migrate-down:
	migrate -path=./migrations -database $$DB_URL down