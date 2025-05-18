package dto

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Fullname string `json:"full_name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type JWTCustomClaims struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	UserRole  string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}
