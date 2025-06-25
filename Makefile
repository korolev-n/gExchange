# Сборка всех контейнеров
build:
	docker compose build

# Запуск всех сервисов
up:
	docker compose up -d

# Остановка
down:
	docker compose down

# Просмотр логов
logs:
	docker compose logs -f

# Перезапуск
restart: down up

# Миграции exchanger
migrate-exchanger:
	docker compose exec exchanger migrate -path=/migrations -database=$$DB_URL up

# Сид данных для exchanger (предполагается, что есть cmd/seed/main.go)
seed-exchanger:
	docker compose exec exchanger go run cmd/seed/main.go

# Миграции wallet
migrate-wallet:
	docker compose exec wallet migrate -path=/migrations -database=$$DB_URL up

# Очистка Docker
clean:
	docker compose down -v
	docker system prune -f

# Проверка здоровья
health-check:
	curl localhost:8080/healthz || true
	curl localhost:8081/healthz || true
