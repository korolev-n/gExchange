# Билд стадии
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY exchanger/go.mod exchanger/go.sum ./ 
# чтобы работал replace
COPY shared/api ./shared/api 
RUN go mod download

COPY exchanger/. .  

RUN CGO_ENABLED=0 GOOS=linux go build -o /exchanger ./cmd/server

# Финальная стадия
FROM alpine:latest

WORKDIR /app

COPY --from=builder /exchanger /app/exchanger
COPY --from=builder /app/migrations /app/migrations

ENV DB_URL=postgres://postgres:postgres@exchanger-db:5432/exchanger?sslmode=disable
ENV SERVER_PORT=8080
ENV LOG_LEVEL=info

EXPOSE 8080 50051

CMD ["/app/exchanger"]