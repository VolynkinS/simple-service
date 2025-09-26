package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"simple-service/pkg/validator"
)

func TestTaskRequestValidation(t *testing.T) {
	tests := []struct {
		name       string
		request    TaskRequest
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Валидный запрос",
			request: TaskRequest{
				Title:       "Валидный заголовок",
				Description: "Валидное описание задачи",
			},
			wantErr: false,
		},
		{
			name: "Пустой title - обязательное поле",
			request: TaskRequest{
				Title:       "",
				Description: "Описание",
			},
			wantErr:    true,
			wantErrMsg: "Field is required for field: Title",
		},
		{
			name: "Слишком длинный title (больше 255 символов)",
			request: TaskRequest{
				Title:       strings.Repeat("a", 256),
				Description: "Описание",
			},
			wantErr:    true,
			wantErrMsg: "Field exceeds maximum length (max: 255 characters) for field: Title",
		},
		{
			name: "Максимально допустимая длина title (255 символов)",
			request: TaskRequest{
				Title:       strings.Repeat("a", 255),
				Description: "Описание",
			},
			wantErr: false,
		},
		{
			name: "Слишком длинное description (больше 1000 символов)",
			request: TaskRequest{
				Title:       "Заголовок",
				Description: strings.Repeat("a", 1001),
			},
			wantErr:    true,
			wantErrMsg: "Field exceeds maximum length (max: 1000 characters) for field: Description",
		},
		{
			name: "Максимально допустимая длина description (1000 символов)",
			request: TaskRequest{
				Title:       "Заголовок",
				Description: strings.Repeat("a", 1000),
			},
			wantErr: false,
		},
		{
			name: "Пустое description допустимо",
			request: TaskRequest{
				Title:       "Заголовок",
				Description: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(context.Background(), tt.request)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
