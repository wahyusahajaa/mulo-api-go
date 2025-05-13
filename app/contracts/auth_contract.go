package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type AuthRepository interface {
	Store(ctx context.Context, input models.RegisterInput) (err error)
	FindUserVerifiedByCode(ctx context.Context, code string) (exists bool, err error)
	FindUserExistsByEmail(ctx context.Context, email string) (exists bool, err error)
	FindUserExistsByUsername(ctx context.Context, username string) (exists bool, err error)
	FindUserByEmail(ctx context.Context, email string) (user *models.User, err error)
	StoreUserVerifyCode(ctx context.Context, userId int, code string) (err error)
	FindUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error)
	UpdateUserVerifiedAt(ctx context.Context, userId int) (err error)
}

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (err error)
	Login(ctx context.Context, req dto.LoginRequest) (token string, err error)
	VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest, userId int) (err error)
	ResendCode(ctx context.Context, userId int) (err error)
}
