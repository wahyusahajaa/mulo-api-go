package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type UserHandler struct {
	svc contracts.UserService
	log *logrus.Logger
}

func NewUserHandler(svc contracts.UserService, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		svc: svc,
		log: log,
	}
}

// @Summary      	List Users
// @Description  	Get paginated list of users
// @Tags         	users
// @Security     	BearerAuth
// @Produce      	json
// @Param        	page     	query    	int  false  "Page number" default(1)
// @Param        	pageSize 	query    	int  false  "Page size" default(10)
// @Success 		200 		{object}	dto.ResponseWithPagination[[]dto.User, dto.Pagination]
// @Failure 		500			{object}	dto.InternalErrorResponse "Internal server error"
// @Router      	/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	users, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "user_handler", "GetUsers", err)
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "user_handler", "GetUsers", err)
	}

	return c.JSON(dto.ResponseWithPagination[[]dto.User, dto.Pagination]{
		Data: users,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

// @Summary      	Get User by ID
// @Description  	Get a user by their ID
// @Tags        	users
// @Security     	BearerAuth
// @Produce      	json
// @Param        	id		path     	int	true  "User ID"
// @Success 		200		{object} 	dto.ResponseWithData[dto.User]
// @Failure 		404 	{object} 	dto.ErrorResponse "User not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router       	/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))

	user, err := h.svc.GetUserById(c.Context(), userId)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "user_handler", "GetUser", err)
	}

	return c.JSON(dto.ResponseWithData[dto.User]{
		Data: user,
	})
}

// UpdateUser		Update an existing user.
// @Summary 		Update user
// @Description 	Updates the user with the specified ID
// @Tags        	users
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "User ID"
// @Param 			user	 body		dto.CreateUserInput true "User object that needs to be updated"
// @Success 		200 	{object} 	dto.ResponseMessage
// @Failure 		400		{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404 	{object} 	dto.ErrorResponse "User not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))
	var req dto.CreateUserInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.Update(c.Context(), req, userId); err != nil {
		return errs.HandleHTTPError(c, h.log, "user_handler", "Update", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Successfully updated user",
	})
}

// UpdateUser		Delete an existing user.
// @Summary 		Delete user
// @Description 	Deletes the user with the specified ID
// @Tags        	users
// @Security     	BearerAuth
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "User ID"
// @Success 		200		{object} 	dto.ResponseMessage
// @Failure 		404 	{object} 	dto.ErrorResponse "User not found"
// @Failure 		500 	{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.Delete(c.Context(), userId); err != nil {
		return errs.HandleHTTPError(c, h.log, "user_handler", "Delete", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted user",
	})
}
