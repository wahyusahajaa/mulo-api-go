package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/pkg/jwt"
)

type AuthMiddleware struct {
	jwtService jwt.JWTService
}

func NewAuthMiddleware(jwtService jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString, err := m.jwtService.ExtractTokenFromHeader(c.Get("Authorization"))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		claims, err := m.jwtService.ParseJWTToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		// Store claims to Locals/Context
		for _, key := range []string{"id", "username", "role"} {
			if value, ok := claims[key]; ok {
				c.Locals(key, value)
			}
		}

		return c.Next()
	}
}
