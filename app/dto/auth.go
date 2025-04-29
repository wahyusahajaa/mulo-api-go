package dto

type RegisterInput struct {
	Fullname string `json:"fullname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyEmailInput struct {
	Code string `json:"code" validate:"required"`
}
