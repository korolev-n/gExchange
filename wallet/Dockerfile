# Build stage
FROM golang:1.24-alpine AS builder

# Создаем структуру проекта внутри контейнера
WORKDIR /workspace
COPY . .

# Переходим в директорию wallet
WORKDIR /workspace/wallet

# Явно указываем replace с абсолютным путем
RUN go mod edit -replace github.com/korolev-n/gExchange/shared/api=/workspace/shared/api

# Проверяем доступность shared/api
RUN ls -la /workspace/shared/api

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /wallet ./cmd/server

# Run stage
FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /wallet .
COPY --from=builder /workspace/wallet/migrations ./migrations

ENV DB_URL=postgres://postgres:postgres@wallet-db:5432/wallet?sslmode=disable
ENV SERVER_PORT=8080