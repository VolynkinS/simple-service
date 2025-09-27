package dto

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// DTO  некоторых компаниях используется такой подход

const (
	FieldBadFormat     = "FIELD_BADFORMAT"
	FieldIncorrect     = "FIELD_INCORRECT"
	ServiceUnavailable = "SERVICE_UNAVAILABLE"
	NotFound           = "NOT_FOUND"
	InternalError      = "Service is currently unavailable. Please try again later."
)

// Swagger DTO structures

// TaskRequest represents the request body for creating a task
// @Description Task creation request
type TaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255" example:"Implement new feature"`
	Description string `json:"description" validate:"max=1000" example:"Develop a new API endpoint for user management"`
} // @name TaskRequest

// TaskResponse represents a task in responses
// @Description Task information
type TaskResponse struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Implement new feature"`
	Description string    `json:"description" example:"Develop a new API endpoint for user management"`
	Status      string    `json:"status" example:"new"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
} // @name TaskResponse

// CreateTaskResponse represents the response after creating a task
// @Description Response after task creation
type CreateTaskResponse struct {
	TaskID int `json:"task_id" example:"1"`
} // @name CreateTaskResponse

// Response represents the standard API response
// @Description Standard API response
type Response struct {
	Status string `json:"status" example:"success"`
	Error  *Error `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
} // @name Response

// SuccessResponse represents a successful API response
// @Description Successful API response
type SuccessResponse struct {
	Status string `json:"status" example:"success"`
	Data   any    `json:"data"`
} // @name SuccessResponse

// ErrorResponse represents an error API response
// @Description Error API response
type ErrorResponse struct {
	Status string `json:"status" example:"error"`
	Error  *Error `json:"error"`
} // @name ErrorResponse

// Error represents error details
// @Description Error details
type Error struct {
	Code string `json:"code" example:"FIELD_INCORRECT"`
	Desc string `json:"desc" example:"Invalid request body"`
} // @name Error

func BadResponseError(ctx *fiber.Ctx, code, desc string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: code,
			Desc: desc,
		},
	})
}

func InternalServerError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: ServiceUnavailable,
			Desc: InternalError,
		},
	})
}

func NotFoundError(ctx *fiber.Ctx, desc string) error {
	return ctx.Status(fiber.StatusNotFound).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: NotFound,
			Desc: desc,
		},
	})
}
