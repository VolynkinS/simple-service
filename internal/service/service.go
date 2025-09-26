package service

import (
	"context"

	"go.uber.org/zap"
)

// Слой бизнес-логики. Тут должна быть основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(ctx context.Context, req TaskRequest) (int, error)
}

// Task - модель задачи для бизнес-логики
type Task struct {
	Title       string
	Description string
}

// Repository - интерфейс для работы с задачами (только в service слое)
type Repository interface {
	CreateTask(ctx context.Context, task Task) (int, error)
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
