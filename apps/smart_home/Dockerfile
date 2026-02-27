# stage сборки приложения
FROM golang:1.20-alpine AS build               # используем официальный golang образ
WORKDIR /src                                   # рабочая директория
COPY . .                                       # копируем исходники в контейнер
RUN apk add --no-cache git                     # устанавливаем git (для go mod)
RUN go mod download                            # скачиваем зависимости
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./... # собираем бинарник статически

# финальный минимальный образ
FROM alpine:3.18                                 # берем маленький alpine для запуска
RUN apk add --no-cache ca-certificates           # добавляем корневые сертификаты
COPY --from=build /app/app /app/app              # копируем бинарник из стадии сборки
EXPOSE 8081                                      # декларируем порт
ENTRYPOINT ["/app/app"]                          # точка входа