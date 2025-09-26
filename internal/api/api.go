package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"

	"simple-service/internal/api/handlers"
	"simple-service/internal/api/middleware"
	"simple-service/internal/service"
)

// Routers - структура для хранения зависимостей роутов
type Routers struct {
	Service service.Service
	Logger  *zap.SugaredLogger
}

// NewRouters - конструктор для настройки API
func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PUT, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	// Группа маршрутов с авторизацией
	apiGroup := app.Group("/v1", middleware.JWTAuthorization(token))

	// Инициализация обработчиков
	taskHandler := handlers.NewTaskHandler(r.Service, r.Logger)

	// Роут для создания задачи
	apiGroup.Post("/create_task", taskHandler.CreateTask)

	return app
}
