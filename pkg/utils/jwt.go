package utils

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
)

type JWTService interface {
	GenerateJWTToken(id int, fullname string, username string, role string, isVerifiedAt bool) (string, error)
	ParseJWTToken(tokenString string) (jwt.MapClaims, error)
	ExtractTokenFromHeader(authHeader string) (string, error)
}

type jwtService struct {
	secret string
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		secret: cfg.JwtSecret,
	}
}

func (j *jwtService) GenerateJWTToken(id int, fullname string, username string, role string, isEmailVerified bool) (string, error) {
	claims := jwt.MapClaims{
		"id":                id,
		"full_name":         fullname,
		"username":          username,
		"role":              role,
		"is_email_verified": isEmailVerified,
		"exp":               time.Now().Add(time.Hour * 24).Unix(),
		// "exp": time.Now().Add(time.Minute * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *jwtService) ParseJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(j.secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil || !token.Valid {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token: "+err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	// Manual cek expiration
	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Token expired")
	}

	return claims, nil
}

func (j *jwtService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header format")
	}

	return parts[1], nil
}
