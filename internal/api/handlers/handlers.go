package handlers

import (
	"encoding/json"
	"simple-service/internal/dto"
	"simple-service/internal/service"
	"simple-service/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TaskHandler struct {
	service service.Service
	log     *zap.SugaredLogger
}

func NewTaskHandler(svc service.Service, logger *zap.SugaredLogger) *TaskHandler {
	return &TaskHandler{
		service: svc,
		log:     logger,
	}
}

func (h *TaskHandler) CreateTask(ctx *fiber.Ctx) error {
	var req service.TaskRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		h.log.Errorw("Invalid request body", "error", err)
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	taskID, err := h.service.CreateTask(ctx.Context(), req)
	if err != nil {
		h.log.Errorw("Failed to create task", "error", err)
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": taskID},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
