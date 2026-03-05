# Этап сборки
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Синхронизируем зависимости и проверяем go.sum
RUN go mod tidy

# Собираем приложение
RUN go build -o smart-temperature-system ./cmd/api

# Этап выполнения
FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN adduser -D -g '' appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/smart-temperature-system .
EXPOSE 8080
CMD ["./smart-temperature-system"]