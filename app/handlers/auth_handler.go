package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type AuthHandler struct {
	svc contracts.AuthService
	log *logrus.Logger
}

func NewAuthHandler(svc contracts.AuthService, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		svc: svc,
		log: log,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.Register(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "Register", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully registered user.",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	token, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "Login", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully logged in.",
		"token":   token,
	})
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var req dto.VerifyEmailRequest
	userId, _ := utils.ParseUserId(c)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.VerifyEmail(c.Context(), req, userId); err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "VerifyEmail", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully user verify",
	})
}

func (h *AuthHandler) ResendCodeEmailVerification(c *fiber.Ctx) error {
	userId := utils.GetUserId(c.Context())

	if err := h.svc.ResendCode(c.Context(), userId); err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "ResendCodeEmailVerification", err)
	}

	return c.JSON(fiber.Map{
		"message": "Successfully resend code.",
	})
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	fullname := c.Locals("full_name")

	return c.JSON(fiber.Map{
		"full_name": fullname,
	})
}
