# Используем официальный образ Go
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Убеждаемся, что все модули инициализированы
RUN go mod tidy

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Используем минимальный образ для продакшена
FROM alpine:latest

# Устанавливаем CA сертификаты и timezone
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -g 1001 appgroup && adduser -u 1001 -G appgroup -s /bin/sh -D appuser

WORKDIR /root/

# Копируем собранное приложение из builder stage
COPY --from=builder /app/main .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Меняем владельца файлов
RUN chown -R appuser:appgroup /root

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
