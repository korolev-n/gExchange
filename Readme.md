**Цель проекта:** улучшение навыков разработки микросервисов (gRPC + HTTP REST).

## Технологии

- **Язык**: Go (v1.24)
- **API**: HTTP REST, gRPC
- **База данных**: PostgreSQL
- **Контейнеризация**: Docker, Docker Compose
- **Кэш**: В памяти + singleflight
- **JWT**: golang-jwt/jwt
- **Логирование**: slog

## `exchanger` — сервис обмена валют

- Отвечает за хранение и предоставление курсов валют.
- Поддерживает два интерфейса:
    - HTTP (`GET /rates`)
    - gRPC (`GetRates`)
- Хранит курсы в PostgreSQL, таблица `exchange_rates`.
- Отдаёт последние курсы (по максимальной дате `set_date`).
- gRPC-сервер использует protobuf-интерфейс `ExchangerService`.

## `wallet` — кошелёк пользователя

- Отвечает за регистрацию, логин, пополнение, снятие и обмен валют.
- REST API с авторизацией на основе JWT.
- Подключается к `exchanger` по **gRPC** для получения курсов валют.
- Использует кэш (`ExchangeRateCache`) для курсов, чтобы избежать лишних gRPC-запросов.
- Поддерживает операции:
    - /register, /login
    - /balance
    - /wallet/deposit, /wallet/withdraw
    - /exchange (через gRPC -> exchanger)

## Взаимодействие сервисов

```
client --> wallet (REST) --> exchanger (gRPC)
                   ↑
               cache TTL
```

- Курсы валют кэшируются.
- Если кэш устарел, `wallet` вызывает `GetRates()` у `exchanger`.

## Конфигурация `.env`

### wallet/.env

- `DB_URL` — строка подключения к PostgreSQL
- `JWT_SECRET` — секрет для подписи JWT
- `JWT_EXPIRATION_HOURS` — срок действия токена
- `SERVER_PORT` — порт HTTP API
- `LOG_LEVEL` - детализация логирования (debug, info, warn)

```
// .env.example

DB_URL=postgres://user:pa55word@localhost:5432/wallet?sslmode=disable
SERVER_PORT=8080
LOG_LEVEL=debug
JWT_EXPIRATION_HOURS=24
JWT_SECRET=secret_key
```
### exchanger/.env

- `DB_URL` — строка подключения к PostgreSQL
- `SERVER_PORT` — порт HTTP API
- `LOG_LEVEL` - детализация логирования (debug, info, warn)

## Инструкция запуска через Docker

### Выполните миграции для обоих сервисов

```bash
make migrate-exchanger
make migrate-wallet
```

### Добавить курсы валют

```bash
make seed-exchanger
```

### Собрать и запустить контейнеры

```bash
make build
make up
```

## Примеры `curl`-запросов

```bash
# Регистрация
curl -X POST localhost:8080/register \
 -H "Content-Type: application/json" \
 -d '{"email": "user@example.com", "password": "secret"}'

# Логин
curl -X POST localhost:8080/login \
 -H "Content-Type: application/json" \
 -d '{"email": "user@example.com", "password": "secret"}'

# Получить баланс
curl -H "Authorization: Bearer <TOKEN>" localhost:8080/balance

# Пополнение
curl -X POST localhost:8080/wallet/deposit \
 -H "Authorization: Bearer <TOKEN>" \
 -H "Content-Type: application/json" \
 -d '{"currency": "USD", "amount": 100}'

# Обмен
curl -X POST localhost:8080/exchange \
 -H "Authorization: Bearer <TOKEN>" \
 -H "Content-Type: application/json" \
 -d '{"from_currency": "USD", "to_currency": "EUR", "amount": 50}'
```
