# Билд стадии
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY wallet/go.mod wallet/go.sum ./
# proto-модуль до загрузки зависимостей
COPY shared/api ./shared/api 
COPY wallet/. .
RUN go mod download

#COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet ./cmd/server

# Финальная стадия
FROM alpine:latest

WORKDIR /app

COPY --from=builder /wallet /app/wallet
COPY --from=builder /app/migrations /app/migrations

ENV DB_URL=postgres://postgres:postgres@wallet-db:5432/wallet?sslmode=disable
ENV SERVER_PORT=8081
ENV LOG_LEVEL=info
ENV JWT_SECRET=secret
ENV JWT_EXPIRATION_HOURS=24
ENV EXCHANGER_GRPC_ADDR=exchanger:50051

EXPOSE 8081

CMD ["/app/wallet"]