package jwt

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
)

type jwtService struct {
	secret string
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		secret: cfg.JwtSecret,
	}
}

func (j *jwtService) GenerateJWTToken(id int, username string, role string) (string, error) {
	claims := jwtlib.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		// "exp": time.Now().Add(time.Minute * 1).Unix(),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) ParseJWTToken(tokenString string) (jwtlib.MapClaims, error) {
	token, err := jwtlib.Parse(tokenString, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(j.secret), nil
	}, jwtlib.WithValidMethods([]string{"HS256"}))

	if err != nil || !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	claims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Authentication token is not valid.")
	}

	// Manual cek expiration
	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Token is expired or no longer valid.")
	}

	return claims, nil
}

func (j *jwtService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "No Authorization header provided. Please include a valid token.")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Authorization header must be in the format: Bearer <token>.")
	}

	return parts[1], nil
}
