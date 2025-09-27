package service

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Слой бизнес-логики. Тут должна быть основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(ctx context.Context, req TaskRequest) (int, error)
	GetTask(ctx context.Context, id int) (*TaskResponse, error)
}

// Task - модель задачи для бизнес-логики
type Task struct {
	Title       string
	Description string
}

// TaskResponse - модель ответа с полной информацией о задаче
type TaskResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Repository - интерфейс для работы с задачами (только в service слое)
type Repository interface {
	CreateTask(ctx context.Context, task Task) (int, error)
	GetTask(ctx context.Context, id int) (*TaskResponse, error)
}

type service struct {
	repo Repository
	log  *zap.SugaredLogger
}

// NewService - конструктор сервиса
func NewService(repository Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repository,
		log:  logger,
	}
}

// CreateTask - бизнес-логика создания задачи
func (s *service) CreateTask(ctx context.Context, req TaskRequest) (int, error) {
	task := req.ToTask()

	taskID, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		s.log.Errorw("Failed to insert task", "error", err)
		return 0, err
	}

	return taskID, nil
}

// GetTask - бизнес-логика получения задачи по ID
func (s *service) GetTask(ctx context.Context, id int) (*TaskResponse, error) {
	task, err := s.repo.GetTask(ctx, id)
	if err != nil {
		s.log.Errorw("Failed to get task", "error", err, "task_id", id)
		return nil, err
	}

	return task, nil
}
