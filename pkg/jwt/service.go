package jwt

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
)

type jwtService struct {
	JwtSecret           []byte
	RefreshSecret       []byte
	AccessTokenExpires  time.Duration
	RefreshTokenExpires time.Duration
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{
		JwtSecret:           []byte(cfg.JwtSecret),
		RefreshSecret:       []byte(cfg.RefreshSecret),
		AccessTokenExpires:  1 * time.Minute,
		RefreshTokenExpires: 7 * 24 * time.Hour,
	}
}

func (j *jwtService) GenerateTokens(id int, username string, role string) (accessToken, refreshToken string, err error) {
	// Generate Access Token
	accessClaims := dto.JWTCustomClaims{
		ID:        id,
		Username:  username,
		UserRole:  role,
		TokenType: "access",
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(j.AccessTokenExpires)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(id),
		},
	}
	accessJwt := jwtlib.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJwt.SignedString(j.JwtSecret)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshClaims := dto.JWTCustomClaims{
		ID:        id,
		Username:  username,
		UserRole:  role,
		TokenType: "refresh",
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(j.RefreshTokenExpires)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			ID:        strconv.Itoa(id),
		},
	}
	refreshJwt := jwtlib.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJwt.SignedString(j.RefreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *jwtService) ParseToken(tokenString, tokenType string) (claims *dto.JWTCustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.JWTCustomClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		if tokenType == "access" {
			return j.JwtSecret, nil
		} else {
			return j.RefreshSecret, nil
		}
	}, jwtlib.WithValidMethods([]string{"HS256"}))
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return nil, errs.NewForbiddenError("Token is expired or no longer valid.")
		}
		if errors.Is(err, jwtlib.ErrSignatureInvalid) {
			return nil, errs.NewForbiddenError("Token signature is invalid.")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*dto.JWTCustomClaims)
	if !ok || !token.Valid {
		return nil, errs.NewForbiddenError("Authentication token is not valid.")
	}

	return claims, nil
}

func (j *jwtService) ParseAccessToken(tokenString string) (claims *dto.JWTCustomClaims, err error) {
	claims, err = j.ParseToken(tokenString, "access")
	if err != nil {
		return nil, err
	}

	// If not access token
	if claims.TokenType != "access" {
		return nil, errs.NewForbiddenError("Invalid token type.")
	}

	return claims, nil
}

func (j *jwtService) ParseRefreshToken(tokenString string) (claims *dto.JWTCustomClaims, err error) {
	claims, err = j.ParseToken(tokenString, "refresh")
	if err != nil {
		return nil, err
	}

	// If not refresh token
	if claims.TokenType != "refresh" {
		return nil, errs.NewForbiddenError("Invalid token type.")
	}

	return claims, nil
}

func (j *jwtService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errs.NewForbiddenError("No Authorization header provided. Please include a valid token.")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errs.NewForbiddenError("Authorization header must be in the format: Bearer <token>.")
	}

	return parts[1], nil
}

func (j *jwtService) AddTokenCookies(c *fiber.Ctx, accessToken string, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    accessToken,
		Expires:  time.Now().Add(j.AccessTokenExpires),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "none",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshTokenExpires),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "none",
	})
}

func (j *jwtService) ClearTokenCookies(c *fiber.Ctx) {
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
}
