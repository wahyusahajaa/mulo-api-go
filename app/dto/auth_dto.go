package dto

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Fullname string `json:"full_name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
} //@name RegisterRequest

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
} //@name LoginRequest

type VerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
} //@name VerifyRequest

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
} //@name ResendVerificationRequest

type JWTCustomClaims struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	UserRole  string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
} //@name JWTCustomClaims

type GithubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
} //@name GithubUser
type GithubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
} //@name GithubEmail

type GithubReq struct {
	Code string `json:"code" validate:"required"`
} //@name GithubReq
