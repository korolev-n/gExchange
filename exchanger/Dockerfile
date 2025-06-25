# Build stage
FROM golang:1.24-alpine AS builder

# Создаем структуру проекта внутри контейнера
WORKDIR /workspace
COPY . .

# Переходим в директорию exchanger и настраиваем модули
WORKDIR /workspace/exchanger

# Явно указываем replace с абсолютным путем
RUN go mod edit -replace github.com/korolev-n/gExchange/shared/api=/workspace/shared/api

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /exchanger ./cmd/server

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /exchanger .
COPY --from=builder /workspace/exchanger/migrations ./migrations
COPY --from=builder /workspace/exchanger/exchanger_shema.yaml .

ENV DB_URL=postgres://postgres:postgres@exchanger-db:5432/exchanger?sslmode=disable
ENV SERVER_PORT=8080

EXPOSE 8080 50051

CMD ["./exchanger"]