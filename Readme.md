**Цель проекта:** улучшение навыков разработки микросервисов (gRPC + HTTP REST).

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
