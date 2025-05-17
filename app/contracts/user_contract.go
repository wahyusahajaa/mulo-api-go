package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type UserRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (users []models.User, err error)
	FindExistsUserByUserID(ctx context.Context, userID int) (exists bool, err error)
	FindUserVerifiedByCode(ctx context.Context, code string) (exists bool, err error)
	FindUserExistsByEmail(ctx context.Context, email string) (exists bool, err error)
	FindUserExistsByUsername(ctx context.Context, username string) (exists bool, err error)
	FindUserByEmail(ctx context.Context, email string) (user *models.User, err error)
	FindUserVerifiedByUserIDAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error)
	FindUserByUserID(ctx context.Context, userID int) (user *models.User, err error)
	Count(ctx context.Context) (total int, err error)
	Update(ctx context.Context, input models.CreateUserInput, userID int) (err error)
	Delete(ctx context.Context, userID int) (err error)
}

type UserService interface {
	GetAll(ctx context.Context, pageSize, offset int) (users []dto.User, err error)
	GetCount(ctx context.Context) (total int, err error)
	GetUserById(ctx context.Context, userID int) (user dto.User, err error)
	Update(ctx context.Context, req dto.CreateUserInput, userID int) (err error)
	Delete(ctx context.Context, userID int) (err error)
}
