package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type AuthHandler struct {
	authService         contracts.AuthService
	jwtService          utils.JWTService
	resendService       utils.ResendService
	verificationService utils.VerificationService
}

func NewAuthHandler(authService contracts.AuthService, jwtService utils.JWTService, resendService utils.ResendService, verificationService utils.VerificationService) *AuthHandler {
	return &AuthHandler{
		authService:         authService,
		jwtService:          jwtService,
		resendService:       resendService,
		verificationService: verificationService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input dto.Credentials

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error on login requrest",
			"errors":  err.Error(),
		})
	}

	// Validate request body
	if errorsMap, err := utils.RequestValidate(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errorsMap,
		})
	}

	user, err := h.authService.GetUserByEmail(c.Context(), input.Email)
	if err != nil {
		log.Printf("failed while get user by email: %v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
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
		log.Printf("SignedString err: %v", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Success login",
		"token":   token,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input dto.RegisterInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error on register request", "errors": err.Error()})
	}

	// Validate request body
	if errorsMap, err := utils.RequestValidate(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errorsMap,
		})
	}

	existsEmail, err := h.authService.CheckUserDuplicateEmail(c.Context(), input.Email)
	if err != nil {
		log.Printf("Check email error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed while register",
		})
	}

	if existsEmail {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	existsUsername, err := h.authService.CheckUserDuplicateUsername(c.Context(), input.Username)
	if err != nil {
		log.Printf("Check username error: %v", err)
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
		log.Printf("failed while generate verification code: %v", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Hash password
	hashPassword := utils.HashPassword(input.Password)
	if err := h.authService.Create(c.Context(), input.Fullname, input.Username, input.Email, hashPassword, code); err != nil {
		log.Printf("Register user error: %v", err.Error())
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

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var input dto.VerifyEmailInput
	idFloat, _ := c.Locals("id").(float64)
	userId := int(idFloat)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error on login requrest",
			"errors":  err.Error(),
		})
	}

	// Validate request body
	if errorsMap, err := utils.RequestValidate(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errorsMap,
		})
	}

	// Find Verify Code by UserId and Code
	exists, err := h.authService.GetUserVerifiedByUserIdAndCode(c.Context(), userId, input.Code)
	if err != nil {
		log.Printf("Get user verified by user_id and code err: %v\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Verify code not found",
		})
	}

	// Update user verified at
	if err := h.authService.UpdateUserVerifiedAt(c.Context(), userId); err != nil {
		log.Printf("update user verified at err: %v\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	data := make(map[string]any)
	data["user_id"] = userId
	data["full_name"] = c.Locals("full_name")

	return c.JSON(fiber.Map{
		"data":    data,
		"message": "Successfully verify email",
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	fullname := c.Locals("full_name")

	return c.JSON(fiber.Map{
		"full_name": fullname,
	})
}
