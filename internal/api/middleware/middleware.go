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

		// Требуем точный префикс "Bearer " с пробелом
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return unauthorizedResponse(c, "Authorization header must start with 'Bearer '")
		}

		// Извлекаем токен после "Bearer "
		token := strings.TrimSpace(authHeader[7:])
		if token == "" {
			return unauthorizedResponse(c, "Authorization token is required")
		}

		// Парсим JWT токен с проверкой алгоритма подписи
		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			// Защита от "alg=none" и algorithm confusion атак
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
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
