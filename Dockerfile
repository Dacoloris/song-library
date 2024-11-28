# Укажите базовый образ
FROM golang:1.23 AS builder

# Установите рабочую директорию
WORKDIR /app

# Скопируйте go.mod и go.sum
COPY go.mod go.sum ./

# Загрузите зависимости
RUN go mod download

# Скопируйте исходный код приложения
COPY . .

# Скомпилируйте приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# Создайте минимальный образ для запуска приложения
FROM alpine:latest

# Установите необходимые пакеты
RUN apk --no-cache add ca-certificates

# Скопируйте скомпилированное приложение из образа builder
COPY --from=builder /app/app .
COPY --from=builder /app/.env .
# Установите переменную окружения
ENV PORT=8080

# Откройте порт
EXPOSE 8080

# Запустите приложение
CMD ["./app"]
