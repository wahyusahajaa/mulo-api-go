package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
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

func (svc *userService) GetAll(ctx context.Context, pageSize, offset int) (users []dto.User, err error) {
	results, err := svc.repo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "user_service", "GetAll", err)
		return nil, err
	}

	users = make([]dto.User, 0, len(results))
	for _, result := range results {
		user := dto.User{
			Id:       result.Id,
			Fullname: result.Fullname,
			Username: result.Username,
			Email:    result.Email,
			Image:    utils.ParseImageToJSON(result.Image),
		}

		users = append(users, user)
	}

	return users, nil
}

func (svc *userService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.Count(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "user_service", "GetCount", err)
		return 0, err
	}

	return
}

func (svc *userService) GetUserById(ctx context.Context, userId int) (user dto.User, err error) {
	result, err := svc.repo.FindUserById(ctx, userId)
	if err != nil {
		utils.LogError(svc.log, ctx, "user_service", "GetUserById", err)
		return user, err
	}
	if result == nil {
		notFoundErr := errs.NewNotFoundError("User", "id", userId)
		utils.LogWarn(svc.log, ctx, "user_service", "GetUserById", notFoundErr)
		return user, notFoundErr
	}

	user = dto.User{
		Id:       result.Id,
		Fullname: result.Fullname,
		Username: result.Username,
		Email:    result.Email,
		Image:    utils.ParseImageToJSON(result.Image),
	}

	return user, nil
}

func (svc *userService) Update(ctx context.Context, req dto.CreateUserInput, userId int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return errs.NewBadRequestError("validation failed", map[string]string{"image": "Invalid image object"})
	}

	exists, err := svc.repo.FindExistsUserById(ctx, userId)
	if err != nil {
		utils.LogError(svc.log, ctx, "user_service", "Update", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("User", "id", userId)
		utils.LogWarn(svc.log, ctx, "user_service", "Update", notFoundErr)
		return notFoundErr
	}

	input := models.CreateUserInput{
		Fullname: req.Fullname,
		Image:    imgByte,
	}

	if err := svc.repo.Update(ctx, input, userId); err != nil {
		utils.LogError(svc.log, ctx, "user_service", "Update", err)
		return err
	}

	return nil
}

func (svc *userService) Delete(ctx context.Context, userId int) (err error) {
	exists, err := svc.repo.FindExistsUserById(ctx, userId)
	if err != nil {
		utils.LogError(svc.log, ctx, "user_service", "Delete", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("User", "id", userId)
		utils.LogWarn(svc.log, ctx, "user_service", "Delete", notFoundErr)
		return notFoundErr
	}

	if err := svc.repo.Delete(ctx, userId); err != nil {
		utils.LogError(svc.log, ctx, "user_service", "Delete", err)
		return err
	}

	return nil
}
