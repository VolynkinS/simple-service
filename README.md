## **Описание проекта**

Simple Service

## API Documentation

Проект использует автоматическую генерацию Swagger документации из аннотаций в коде.

### Просмотр документации

1. **Онлайн**: Запустите сервис и перейдите по адресу `http://localhost:8080/swagger/`
2. **Файл**: Откройте `docs/swagger.yaml` в Swagger Editor

### Генерация документации

```bash
# Установка swag CLI (если не установлен)
make swagger-install

# Генерация документации из аннотаций
make swagger-gen

# Или напрямую
swag init -g cmd/main.go -o docs
```

### Структура документации

- `docs/swagger.yaml` - YAML спецификация
- `docs/swagger.json` - JSON спецификация
- `docs/docs.go` - Go код для встраивания

### Добавление новых эндпоинтов

1. Добавь аннотации к handler функции:
```go
// @Summary Краткое описание
// @Description Подробное описание
// @Tags тег
// @Accept json
// @Produce json
// @Param request body dto.RequestType true "Описание параметра"
// @Success 200 {object} dto.ResponseType
// @Failure 400 {object} dto.ErrorResponse
// @Router /v1/endpoint [post]
func (h *Handler) Method(ctx *fiber.Ctx) error {
    // код handler'а
}
```

2. Создай DTO структуры в `internal/dto/`:
```go
// RequestType описание запроса
// @Description Описание для Swagger
type RequestType struct {
    Field string `json:"field" example:"пример"`
} // @name RequestType
```

3. Перегенерируй документацию:
```bash
make swagger-gen
```

### Команды Makefile

```bash
make swagger-gen     # Генерация документации
make swagger-serve   # Информация о URL Swagger UI
make build          # Сборка приложения
make run            # Запуск приложения
make all            # Полная сборка с документацией
```

## Simple Service – это REST API-сервис, написанный на Go с использованием фреймворка Fiber и PostgreSQL. Сервис предоставляет базовый функционал для управления задачами.

Реализовано:

- Создание задач через API
- Валидация входных данных
- Логирование с использованием `zap`
- Хранение данных в PostgreSQL
- Подключение через `pgxpool` для эффективного управления соединениями

---

## **1️⃣ Подготовка окружения**

### **1.1 Установка зависимостей**

Перед запуском убедитесь, что у вас установлены:

- Git
- Go
- Docker
- DataGrip
- Postman для тестирования API

### **1.2 Клонирование репозитория**

```
git clone https://github.com/yourusername/simple-service.git
cd simple-service
```

---

## **2️⃣ Запуск PostgreSQL в Docker**

### **2.1 Запуск контейнера**

Создайте и запустите контейнер с PostgreSQL:

```
docker run --name postgres-db -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=simple_service -p 5432:5432 -d postgres:latest

```

**Параметры:**

- `POSTGRES_USER=admin` – имя пользователя БД
- `POSTGRES_PASSWORD=admin` – пароль
- `POSTGRES_DB=simple_service` – название базы
- `-p 5432:5432` – проброс порта

### **2.2 Проверка работы БД**

Подключитесь к PostgreSQL с помощью программы (DataGrip) и создайте там таблицу

---

## **3️⃣ Настройка проекта**

### **3.1 `local.env` файл**

У вас он уже есть в проекте. В `.env` файлах лежат переменные окружения, которые нельзя хранить в коде для безопасности, и конфигурация сервиса (пример):

```
POSTGRES_USER=admin
POSTGRES_PASSWORD=admin
POSTGRES_DB=simple_service
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
REST_LISTEN_ADDRESS=:8080
REST_TOKEN=your_secret_token
```

### **3.2 Применение миграций**

Создайте таблицу `tasks` в базе данных:

```
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

```
---

## **4️⃣ Запуск сервиса**

### **4.1 Локальный запуск**
Таким способом мы **не** запускаем проекты во время локальной разработки:
```
go run cmd/main.go
```
Всегда запускайте в IDE в **Debug** или в обычном режимах. Описано в pdf файле в задании на kaiton.

Сервис будет доступен по адресу `http://localhost:8080`, если в `.env` файле вы указали PORT=:8080.

---

## **5️⃣ Тестирование API**

### **5.1 Создание задачи**

**Запрос:**

```
POST http://localhost:8080/v1/create_task
Content-Type: application/json
Authorization: Bearer your_secret_token

```

```
{
  "title": "New Feature",
  "description": "Develop new API endpoint"
}

```

**Ответ:**

```
{
  "status": "success",
  "data": {
    "task_id": 1
  }
}

```

---

## **6️⃣ Остановка и удаление контейнера**

```
docker stop postgres-db && docker rm postgres-db

```
---

## **Дополнительная информация**

- Файл `docs/openapi.yaml` содержит документацию API в формате OpenAPI 3.0
- Логирование ведётся через `zap.Logger`
- Переменные окружения загружаются через `envconfig`
- Соединение с PostgreSQL осуществляется через `pgxpool`

Сервис готов к работе.
