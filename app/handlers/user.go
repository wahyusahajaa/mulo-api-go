package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
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

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)

	users, err := h.svc.GetAll(c.Context(), pageSize, offset)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "user_handler", "GetUsers", err)
	}

	total, err := h.svc.GetCount(c.Context())
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "user_handler", "GetUsers", err)
	}

	return c.JSON(fiber.Map{
		"data": users,
		"pagination": dto.Pagination{
			PageSize: pageSize,
			Page:     page,
			Total:    total,
		},
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))

	user, err := h.svc.GetUserById(c.Context(), userId)
	if err != nil {
		return utils.HandleHTTPError(c, h.log, "user_handler", "GetUser", err)
	}

	return c.JSON(fiber.Map{
		"data": user,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))
	var req dto.CreateUserInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.Update(c.Context(), req, userId); err != nil {
		return utils.HandleHTTPError(c, h.log, "user_handler", "Update", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully updated user",
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.Delete(c.Context(), userId); err != nil {
		return utils.HandleHTTPError(c, h.log, "user_handler", "Delete", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully deleted user",
	})
}
