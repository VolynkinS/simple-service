package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthorization - middleware для проверки JWT токена
func JWTAuthorization(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return unauthorizedResponse(c, "Authorization header is required")
		}

		if !strings.HasPrefix(authHeader, "Bearer") {
			return unauthorizedResponse(c, "Authorization header must start with 'Bearer '")
		}

		// Обрабатываем случаи "Bearer" и "Bearer "
		if authHeader == "Bearer" || authHeader == "Bearer " {
			return unauthorizedResponse(c, "Authorization token is required")
		}

		// Извлекаем токен после "Bearer "
		var token string
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimSpace(authHeader[7:])
		} else {
			token = strings.TrimSpace(authHeader[6:])
		}

		if token == "" {
			return unauthorizedResponse(c, "Authorization token is required")
		}

		// Парсим JWT токен
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			return unauthorizedResponse(c, "Invalid authorization token")
		}

		if !parsedToken.Valid {
			return unauthorizedResponse(c, "Invalid authorization token")
		}

		// Сохраняем claims в контекст
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			c.Locals("user", claims)
		}

		return c.Next()
	}
}

func unauthorizedResponse(c *fiber.Ctx, desc string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": "error",
		"error":  fiber.Map{"code": "UNAUTHORIZED", "desc": desc},
	})
}
