package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type UserHandler struct {
	svc contracts.UserService
}

func NewUserHandler(svc contracts.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, pageSize, offset := utils.GetPaginationParam(c)
	users, err := h.svc.GetAll(c.Context(), pageSize, offset)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	total, err := h.svc.GetCount(c.Context())

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var pagination dto.Pagination
	pagination.Total = total
	pagination.PageSize = pageSize
	pagination.Page = page

	results := []dto.User{}

	for _, v := range users {
		result := dto.User{}
		result.Id = v.Id
		result.Fullname = v.Fullname
		result.Username = v.Username
		result.Email = v.Email
		result.Image = utils.ParseImageToJSON(v.Image)
		results = append(results, result)
	}

	return c.JSON(fiber.Map{
		"data":       results,
		"pagination": pagination,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))
	user, err := h.svc.GetUserById(c.Context(), userId)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var result *dto.User

	if user != nil {
		result = &dto.User{
			Id:       user.Id,
			Fullname: user.Fullname,
			Username: user.Username,
			Email:    user.Email,
			Image:    utils.ParseImageToJSON(user.Image),
		}
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))
	var input dto.UserUpdateInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	imgByte, err := utils.ParseImageToByte(input.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid image object",
		})
	}

	user, err := h.svc.GetUserById(c.Context(), userId)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := h.svc.Update(c.Context(), input.Fullname, imgByte, userId); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully update user",
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))

	user, err := h.svc.GetUserById(c.Context(), userId)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := h.svc.Delete(c.Context(), userId); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully delete user",
	})
}
