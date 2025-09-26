# CLAUDE.md

Этот файл предоставляет руководство для Claude Code (claude.ai/code) при работе с кодом в данном репозитории.

## Обзор проекта

Go REST API сервис использующий фреймворк Fiber и PostgreSQL для управления задачами. Сервис следует чистой архитектуре с внедрением зависимостей и структурированным логированием.

## Ключевые команды

### Разработка
- **Запуск сервиса**: Используйте режим Debug в IDE (не командную строку `go run`)
- **Сборка**: `go build -o bin/simple-service cmd/main.go`
- **Тестирование всех пакетов**: `go test ./...`
- **Тестирование конкретного пакета**: `go test ./internal/service -v`
- **Тестирование с покрытием**: `go test -cover ./...`

### Настройка базы данных
```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

Настройка PostgreSQL через Docker:
```bash
docker run --name postgres-db -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=simple_service -p 5432:5432 -d postgres:latest
```

## Архитектура

**Поток зависимостей**: `main.go` → `api` → `service` → `repo` → PostgreSQL

**Ключевые слои**:
- `cmd/main.go` - Точка входа с внедрением зависимостей и graceful shutdown
- `internal/api/` - Fiber HTTP обработчики с CORS и auth middleware
- `internal/service/` - Бизнес-логика с валидацией
- `internal/repo/` - Доступ к данным PostgreSQL через pgxpool
- `internal/config/` - Конфигурация окружения через envconfig

**Конфигурация**: Использует файл `local.env`, загружаемый godotenv и парсится envconfig со структурными тегами.

**Аутентификация**: JWT middleware защищает все маршруты `/v1/*` с проверкой подписи и срока действия токена.

**Валидация**:
- `title` - обязательное поле, 1-255 символов
- `description` - необязательное поле, максимум 1000 символов

**Тестирование**:
- Комплексное тестирование с testify/mock для репозитория
- JWT middleware тесты с валидными/невалидными токенами
- Валидация входных данных с граничными значениями

## Переменные окружения

Обязательные в `local.env`:
- `PORT` - Адрес прослушивания (`:8080`)
- `TOKEN` - Секретный ключ для подписи JWT токенов
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` - Подключение к PostgreSQL
- `WRITE_TIMEOUT`, `SERVER_NAME` - Конфигурация HTTP сервера

## Структура API

Единственная точка входа: `POST /v1/create_task`
- Требует заголовок `Authorization: Bearer <jwt_token>`
- Вход: `{"title": "string (1-255)", "description": "string (макс 1000, необязательно)"}`
- Выход: `{"status": "success", "data": {"task_id": int}}`

## Безопасность

- JWT авторизация с проверкой подписи и истечения срока
- Валидация длины входных данных для предотвращения атак
- Structured logging для аудита безопасности