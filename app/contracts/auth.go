package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type AuthRepository interface {
	Store(ctx context.Context, fullname, username, email, password, code string) (err error)
	FindUserVerifiedByCode(ctx context.Context, code string) (exists bool, err error)
	FindUserDuplicateEmail(ctx context.Context, email string) (exists bool, err error)
	FindUserDuplicateUsername(ctx context.Context, username string) (exists bool, err error)
	FindUserByEmail(ctx context.Context, email string) (user *models.User, err error)
	StoreUserVerifyCode(ctx context.Context, userId int, code string) (err error)
	FindUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error)
	UpdateUserVerifiedAt(ctx context.Context, userId int) (err error)
}

type AuthService interface {
	Create(ctx context.Context, fullname, username, email, password, code string) (err error)
	CheckVerificationCode(ctx context.Context, code string) (exists bool, err error)
	CheckUserDuplicateEmail(ctx context.Context, email string) (exists bool, err error)
	CheckUserDuplicateUsername(ctx context.Context, username string) (exists bool, err error)
	GetUserByEmail(ctx context.Context, email string) (user *models.User, err error)
	CreateUserVerifyCode(ctx context.Context, userId int, code string) (err error)
	GetUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error)
	UpdateUserVerifiedAt(ctx context.Context, userId int) (err error)
}
