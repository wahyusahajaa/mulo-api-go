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

		if !claims["is_email_verified"].(bool) && c.Path() != "/v1/auth/verify-email" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Email not verified yet.",
			})
		}

		// Store claims to Locals
		for _, key := range []string{"id", "full_name", "username", "role", "is_email_verified"} {
			if value, ok := claims[key]; ok {
				c.Locals(key, value)
			}
		}

		return c.Next()
	}
}
