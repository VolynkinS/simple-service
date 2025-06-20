package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"simple-service/internal/dto"
	"simple-service/internal/repo"
	"simple-service/pkg/validator"
	"strconv"
	"time"
)

// Слой бизнес-логики. Тут должна быть основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTask(ctx *fiber.Ctx) error
}

type service struct {
	repo repo.Repository
	log  *zap.SugaredLogger
}

// NewService - конструктор сервиса
func NewService(repo repo.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

// CreateTask - обработчик запроса на создание задачи
func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	// Десериализация JSON-запроса
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Вставка задачи в БД через репозиторий
	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
	}
	taskID, err := s.repo.CreateTask(ctx.Context(), task)
	if err != nil {
		s.log.Error("Failed to insert task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": taskID},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetTask - обработчик запроса на получение задачи по id
func (s *service) GetTask(ctx *fiber.Ctx) error {
	// Получение taskID из адреса запроса
	taskIDstr := ctx.Params("id")
	taskID, err := strconv.Atoi(taskIDstr)
	if err != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "invalid task ID")
	}

	// Получение task из репозитория
	task, err := s.repo.GetTask(ctx.Context(), taskID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data: map[string]string{
			"id":          strconv.Itoa(task.Id),
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"created_at":  task.CreatedAt.Format(time.RFC3339),
			"updated_at":  task.UpdatedAt.Format(time.RFC3339),
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
