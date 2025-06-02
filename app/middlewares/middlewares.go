package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
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
		// tokenString, err := m.jwtService.ExtractTokenFromHeader(c.Get("Authorization"))
		// if err != nil {
		// 	var forbiddenErr *errs.Fobidden
		// 	if errors.As(err, &forbiddenErr) {
		// 		return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
		// 			Message: forbiddenErr.Message,
		// 		})
		// 	}
		// 	return c.Status(fiber.StatusForbidden).JSON(dto.ErrorResponse{
		// 		Message: "Authorization header is missing.",
		// 	})
		// }

		tokenString := c.Cookies("access_token")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.ErrorResponse{
				Message: "Missing access token.",
			})
		}

		claims, err := m.jwtService.ParseAccessToken(tokenString)
		if err != nil {
			var forbiddenErr *errs.Fobidden
			if errors.As(err, &forbiddenErr) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": forbiddenErr.Message})
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
		}

		// Store claims to Locals/Context
		c.Locals("id", claims.ID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.UserRole)
		c.Locals("token_type", claims.TokenType)

		return c.Next()
	}
}
