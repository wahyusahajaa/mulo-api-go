package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type AuthRepository interface {
	Store(ctx context.Context, input models.RegisterInput) (err error)
	StoreUserVerifyCode(ctx context.Context, userId int, code string) (err error)
	UpdateUserVerifiedAt(ctx context.Context, userId int) (err error)
}

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (err error)
	Login(ctx context.Context, req dto.LoginRequest) (token string, err error)
	Verify(ctx context.Context, req dto.VerifyRequest) (err error)
	ResendVerification(ctx context.Context, req dto.ResendVerificationRequest) (err error)
	VerificationStatus(ctx context.Context, email string) (status bool, err error)
	AuthMe(ctx context.Context, userID int) (user dto.User, err error)
}
