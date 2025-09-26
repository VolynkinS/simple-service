package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthorization(t *testing.T) {
	app := fiber.New()
	secretKey := "test-secret-key"

	// Добавляем middleware авторизации
	app.Use(JWTAuthorization(secretKey))

	// Тестовый роут
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	// Создаем валидный JWT токен для тестов
	validToken := createTestJWT(t, secretKey, map[string]interface{}{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	// Создаем истекший JWT токен
	expiredToken := createTestJWT(t, secretKey, map[string]interface{}{
		"user_id": "123",
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})

	// Создаем токен с неправильной подписью
	wrongSignatureToken := createTestJWT(t, "wrong-secret", map[string]interface{}{
		"user_id": "123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Успешная авторизация с валидным JWT",
			authHeader:     "Bearer " + validToken,
			expectedStatus: 200,
			expectedBody:   `{"message":"success"}`,
		},
		{
			name:           "Отсутствует заголовок Authorization",
			authHeader:     "",
			expectedStatus: 401,
			expectedBody:   `{"error":{"code":"UNAUTHORIZED","desc":"Authorization header is required"},"status":"error"}`,
		},
		{
			name:           "Неправильный формат заголовка",
			authHeader:     "Token " + validToken,
			expectedStatus: 401,
			expectedBody:   `{"error":{"code":"UNAUTHORIZED","desc":"Authorization header must start with 'Bearer '"},"status":"error"}`,
		},
		{
			name:           "Пустой токен после Bearer",
			authHeader:     "Bearer ",
			expectedStatus: 401,
			expectedBody:   `{"error":{"code":"UNAUTHORIZED","desc":"Authorization token is required"},"status":"error"}`,
		},
		{
			name:           "Неправильная подпись токена",
			authHeader:     "Bearer " + wrongSignatureToken,
			expectedStatus: 401,
			expectedBody:   `{"error":{"code":"UNAUTHORIZED","desc":"Invalid authorization token"},"status":"error"}`,
		},
		{
			name:           "Истекший токен",
			authHeader:     "Bearer " + expiredToken,
			expectedStatus: 401,
			expectedBody:   `{"error":{"code":"UNAUTHORIZED","desc":"Invalid authorization token"},"status":"error"}`,
		},
		{
			name:           "Некорректный JWT формат",
			authHeader:     "Bearer invalid.jwt.token",
			expectedStatus: 401,
			expectedBody:   `{"error":{"code":"UNAUTHORIZED","desc":"Invalid authorization token"},"status":"error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Проверяем тело ответа для неуспешных случаев
			if tt.expectedStatus != 200 {
				buf := make([]byte, 1024)
				n, _ := resp.Body.Read(buf)
				actualBody := string(buf[:n])
				assert.JSONEq(t, tt.expectedBody, actualBody)
			}
		})
	}
}

// createTestJWT создает JWT токен для тестов
func createTestJWT(t *testing.T, secretKey string, claims map[string]interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(secretKey))
	assert.NoError(t, err)
	return tokenString
}
