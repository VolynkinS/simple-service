package handlers

import (
	"encoding/json"
	"simple-service/internal/dto"
	"simple-service/internal/service"
	"simple-service/pkg/validator"
	"strconv"

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

// CreateTask creates a new task
// @Summary Create a new task
// @Description Creates a new task in the system
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body dto.TaskRequest true "Task data"
// @Success 200 {object} dto.SuccessResponse{data=dto.CreateTaskResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/create_task [post]
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

// GetTask retrieves a task by ID
// @Summary Get task by ID
// @Description Retrieves a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.TaskResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/tasks/{id} [get]
func (h *TaskHandler) GetTask(ctx *fiber.Ctx) error {
	// Получаем ID из параметров URL
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Errorw("Invalid task ID", "error", err, "id", idStr)
		return dto.BadResponseError(ctx, dto.FieldIncorrect, "Invalid task ID")
	}

	// Получаем задачу из сервиса
	task, err := h.service.GetTask(ctx.Context(), id)
	if err != nil {
		h.log.Errorw("Failed to get task", "error", err, "task_id", id)
		if err.Error() == "task not found" {
			return dto.NotFoundError(ctx, "Task not found")
		}
		return dto.InternalServerError(ctx)
	}

	response := dto.Response{
		Status: "success",
		Data:   task,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
