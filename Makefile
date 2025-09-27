.PHONY: swagger-gen swagger-install swagger-serve build run

# Установка swag CLI
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

# Генерация Swagger документации
swagger-gen:
	swag init -g cmd/main.go -o docs

# Запуск Swagger UI локально
swagger-serve:
	@echo "Swagger UI will be available at http://localhost:8080/swagger/"

# Сборка приложения
build:
	go build -o bin/simple-service cmd/main.go

# Запуск приложения
run:
	go run cmd/main.go

# Обновление зависимостей
deps:
	go mod tidy
	go mod download

# Форматирование кода
fmt:
	go fmt ./...

# Линтинг
lint:
	golangci-lint run

# Тесты
test:
	go test -v ./... -cover

# Полная сборка с документацией
all: deps swagger-gen build

# Очистка
clean:
	rm -rf bin/
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml

# Справка
help:
	@echo "Available targets:"
	@echo "  swagger-install - Install swag CLI"
	@echo "  swagger-gen     - Generate Swagger documentation"
	@echo "  swagger-serve   - Info about Swagger UI URL"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  deps           - Update dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  test           - Run tests"
	@echo "  all            - Full build with docs"
	@echo "  clean          - Clean build artifacts"
	@echo "  help           - Show this help"
