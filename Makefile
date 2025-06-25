.PHONY: up down

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

exchanger-logs:
	docker compose logs -f exchanger

wallet-logs:
	docker compose logs -f wallet

restart: down up

clean:
	docker compose down -v

ps:
	docker compose ps