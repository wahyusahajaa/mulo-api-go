package services

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type authService struct {
	repo contracts.AuthRepository
}

func NewAuthService(repo contracts.AuthRepository) contracts.AuthService {
	return &authService{
		repo: repo,
	}
}

func (svc *authService) Create(ctx context.Context, fullname, username, email, password, code string) (err error) {
	return svc.repo.Store(ctx, fullname, username, email, password, code)
}

func (svc *authService) CheckVerificationCode(ctx context.Context, code string) (exists bool, err error) {
	return svc.repo.FindUserVerifiedByCode(ctx, code)
}

func (svc *authService) CheckUserDuplicateEmail(ctx context.Context, email string) (exists bool, err error) {
	return svc.repo.FindUserDuplicateEmail(ctx, email)
}

func (svc *authService) CheckUserDuplicateUsername(ctx context.Context, username string) (exists bool, err error) {
	return svc.repo.FindUserDuplicateUsername(ctx, username)
}

func (svc *authService) GetUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	return svc.repo.FindUserByEmail(ctx, email)
}

func (svc *authService) CreateVerifyCode(ctx context.Context, userId int, code string) (err error) {
	return svc.repo.StoreVerifyCode(ctx, userId, code)
}

func (svc *authService) GetUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (exists bool, err error) {
	return svc.repo.FindUserVerifiedByUserIdAndCode(ctx, userId, code)
}

func (svc *authService) UpdateUserVerifiedAt(ctx context.Context, userId int) (err error) {
	return svc.repo.UpdateUserVerifiedAt(ctx, userId)
}
