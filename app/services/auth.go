package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type authService struct {
	repo contracts.AuthRepository
	log  *logrus.Logger
}

func NewAuthService(repo contracts.AuthRepository, log *logrus.Logger) contracts.AuthService {
	return &authService{
		repo: repo,
		log:  log,
	}
}

func (svc *authService) Create(ctx context.Context, fullname, username, email, password, code string) (err error) {
	err = svc.repo.Store(ctx, fullname, username, email, password, code)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}

func (svc *authService) CheckVerificationCode(ctx context.Context, code string) (exists bool, err error) {
	exists, err = svc.repo.FindUserVerifiedByCode(ctx, code)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}

func (svc *authService) CheckUserDuplicateEmail(ctx context.Context, email string) (exists bool, err error) {
	exists, err = svc.repo.FindUserDuplicateEmail(ctx, email)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}

func (svc *authService) CheckUserDuplicateUsername(ctx context.Context, username string) (exists bool, err error) {
	exists, err = svc.repo.FindUserDuplicateUsername(ctx, username)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}

func (svc *authService) GetUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	user, err = svc.repo.FindUserByEmail(ctx, email)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}

func (svc *authService) CreateUserVerifyCode(ctx context.Context, userId int, code string) (err error) {
	err = svc.repo.StoreUserVerifyCode(ctx, userId, code)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}

func (svc *authService) GetUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error) {
	userVerified, err = svc.repo.FindUserVerifiedByUserIdAndCode(ctx, userId, code)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}

func (svc *authService) UpdateUserVerifiedAt(ctx context.Context, userId int) (err error) {
	err = svc.repo.UpdateUserVerifiedAt(ctx, userId)

	if err != nil {
		svc.log.WithError(err).Error("error in auth service")
	}

	return
}
