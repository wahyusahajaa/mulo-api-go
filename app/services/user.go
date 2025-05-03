package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type userService struct {
	repo contracts.UserRepository
	log  *logrus.Logger
}

func NewUserService(repo contracts.UserRepository, log *logrus.Logger) contracts.UserService {
	return &userService{
		repo: repo,
		log:  log,
	}
}

func (svc *userService) GetAll(ctx context.Context, pageSize, offset int) (users []models.User, err error) {
	users, err = svc.repo.FindAll(ctx, pageSize, offset)

	if err != nil {
		svc.log.WithError(err).Error("error in user service")
		return
	}

	return
}

func (svc *userService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.Count(ctx)

	if err != nil {
		svc.log.WithError(err).Error("error in user service")
		return 0, err
	}

	return
}

func (svc *userService) GetUserById(ctx context.Context, userId int) (user *models.User, err error) {
	user, err = svc.repo.FindUserById(ctx, userId)

	if err != nil {
		svc.log.WithError(err).Error("error in user service")
		return nil, err
	}

	return
}

func (svc *userService) Update(ctx context.Context, fullname string, image []byte, userId int) (err error) {
	err = svc.repo.Update(ctx, fullname, image, userId)

	if err != nil {
		svc.log.WithError(err).Error("error in user service")
		return err
	}

	return nil
}

func (svc *userService) Delete(ctx context.Context, userId int) (err error) {
	err = svc.repo.Delete(ctx, userId)

	if err != nil {
		svc.log.WithError(err).Error("error in user service")
		return err
	}

	return nil
}
