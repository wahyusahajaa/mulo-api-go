package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type AuthHandler struct {
	authService         contracts.AuthService
	jwtService          utils.JWTService
	resendService       utils.ResendService
	verificationService utils.VerificationService
	log                 *logrus.Logger
}

func NewAuthHandler(authService contracts.AuthService, jwtService utils.JWTService, resendService utils.ResendService, verificationService utils.VerificationService, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		jwtService:          jwtService,
		resendService:       resendService,
		verificationService: verificationService,
		log:                 log,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error on register request", "errors": err.Error()})
	}

	// Validate request body
	if errorsMap, err := utils.RequestValidate(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Failed",
			"errors":  errorsMap,
		})
	}

	existsEmail, err := h.authService.CheckUserDuplicateEmail(c.Context(), input.Email)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if existsEmail {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	existsUsername, err := h.authService.CheckUserDuplicateUsername(c.Context(), input.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed while register",
		})
	}

	if existsUsername {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// Generate User Verification Code
	code, err := h.verificationService.GenerateVerificationCode(c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Hash password
	hashPassword := utils.HashPassword(input.Password)
	if err := h.authService.Create(c.Context(), input.Fullname, input.Username, input.Email, hashPassword, code); err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"message": "Failed while register",
		})
	}

	// Send Verification code
	// go h.resendService.SendEmailVerificationCode("wahyusahaja.official@gmail.com", code)

	return c.JSON(fiber.Map{
		"message": "Register successfully",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input dto.Credentials

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request body
	if errorsMap, err := utils.RequestValidate(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Failed",
			"errors":  errorsMap,
		})
	}

	user, err := h.authService.GetUserByEmail(c.Context(), input.Email)
	if err != nil {
		h.log.WithError(err).Error("failed to get user by email")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// Check Password
	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	// Generate token
	token, err := h.jwtService.GenerateJWTToken(user.Id, user.Fullname, user.Username, user.Role, user.EmailVerifiedAt.Valid)

	if err != nil {
		h.log.WithError(err).Error("failed to generate token")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Success login",
		"token":   token,
	})
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var input dto.VerifyEmailInput
	userId, _ := utils.ParseUserId(c)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate request body
	if errorsMap, err := utils.RequestValidate(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation Failed",
			"errors":  errorsMap,
		})
	}

	// Find Verify Code by UserId and Code
	userVerified, err := h.authService.GetUserVerifiedByUserIdAndCode(c.Context(), userId, input.Code)
	if err != nil {
		h.log.WithError(err).Error("failed to get user verified by user id and code")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if userVerified == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User verify code not found",
		})
	}

	if userVerified.ExpiredAt.Valid && userVerified.ExpiredAt.Time.Before(time.Now()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Verify code is expired",
		})
	}

	// Update user verified at
	if err := h.authService.UpdateUserVerifiedAt(c.Context(), userId); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Success verify email",
	})
}

func (h *AuthHandler) ResendCodeEmailVerification(c *fiber.Ctx) error {
	userId, _ := utils.ParseUserId(c)

	code, err := h.verificationService.GenerateVerificationCode(c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := h.authService.CreateUserVerifyCode(c.Context(), userId, code); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Success Resend code verification",
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	fullname := c.Locals("full_name")

	return c.JSON(fiber.Map{
		"full_name": fullname,
	})
}
