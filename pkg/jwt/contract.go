package jwt

import (
	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateJWTToken(id int, username string, role string) (string, error)
	ParseJWTToken(tokenString string) (jwtlib.MapClaims, error)
	ExtractTokenFromHeader(authHeader string) (string, error)
}
