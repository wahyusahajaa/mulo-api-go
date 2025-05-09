package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type UserRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (users []models.User, err error)
	FindExistsUserById(ctx context.Context, userId int) (exists bool, err error)
	Count(ctx context.Context) (total int, err error)
	FindUserById(ctx context.Context, userId int) (user *models.User, err error)
	Update(ctx context.Context, input models.CreateUserInput, userId int) (err error)
	Delete(ctx context.Context, userId int) (err error)
}

type UserService interface {
	GetAll(ctx context.Context, pageSize, offset int) (users []dto.User, err error)
	GetCount(ctx context.Context) (total int, err error)
	GetUserById(ctx context.Context, userId int) (user dto.User, err error)
	Update(ctx context.Context, req dto.CreateUserInput, userId int) (err error)
	Delete(ctx context.Context, userId int) (err error)
}
