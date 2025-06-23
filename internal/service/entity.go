package service

// TaskRequest - структура, представляющая тело запроса на создание задачи
type TaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}
