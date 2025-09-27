package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
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

	// Swagger UI (без авторизации)
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Группа маршрутов с авторизацией
	apiGroup := app.Group("/v1", middleware.JWTAuthorization(token))

	// Инициализация обработчиков
	taskHandler := handlers.NewTaskHandler(r.Service, r.Logger)

	// Роуты для задач
	apiGroup.Post("/create_task", taskHandler.CreateTask)
	apiGroup.Get("/tasks/:id", taskHandler.GetTask)

	return app
}
