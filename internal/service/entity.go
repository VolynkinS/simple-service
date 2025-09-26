package service

// TaskRequest - структура, представляющая тело запроса
type TaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// ToTask - конвертирует TaskRequest в Task
func (tr TaskRequest) ToTask() Task {
	return Task(tr)
}
