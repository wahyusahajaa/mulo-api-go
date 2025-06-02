package jwt

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
)

type JWTService interface {
	GenerateTokens(id int, username string, role string) (accessToken, refreshToken string, err error)
	ParseToken(token, tokenType string) (claims *dto.JWTCustomClaims, err error)
	ParseAccessToken(tokenString string) (claims *dto.JWTCustomClaims, err error)
	ParseRefreshToken(tokenString string) (claims *dto.JWTCustomClaims, err error)
	ExtractTokenFromHeader(authHeader string) (string, error)
	AddTokenCookies(c *fiber.Ctx, accessToken, refreshToken string)
	ClearTokenCookies(c *fiber.Ctx)
}
