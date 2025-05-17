package handlers

import (
	"net/mail"

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

// Register			User registration
// @Summary 		Register a new user
// @Description 	Create a new user account and sends a verification email.
// @Tags        	auth
// @Accept 			json
// @Produce 		json
// @Param 			register	 body	dto.RegisterRequest true "register object that needs to be created"
// @Success 		201 		{object} 	dto.ResponseMessage
// @Failure 		400			{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		409			{object} 	dto.ValidationErrorResponse "Username or email already exists."
// @Failure 		500 		{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseError{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.Register(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "Register", err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ResponseMessage{
		Message: "User registered successfully.",
	})
}

// Login			Login user
// @Summary 		Login user
// @Description 	Authenticates a user and returns a JWT token if successful.
// @Tags        	auth
// @Accept 			json
// @Produce 		json
// @Param 			login	 	body		dto.LoginRequest true "login object that needs to be created"
// @Success 		201 		{object} 	dto.ResponseMessage
// @Failure 		400			{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		403			{object} 	dto.ValidationErrorResponse "Account not activated"
// @Failure 		404			{object} 	dto.ValidationErrorResponse "Invalid email or password"
// @Failure 		500 		{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseError{
			Message: "Invalid body request.",
		})
	}

	token, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "Login", err)
	}

	return c.JSON(dto.ResponseWithToken[string, string]{
		Message: "Successfully logged in.",
		Token:   token,
	})
}

// Verify			Verify user email
// @Summary 		Verify user email
// @Description 	Verifies the user's email address using a verification code.
// @Tags        	auth
// @Accept 			json
// @Produce 		json
// @Param 			verify	 	body		dto.VerifyRequest true "verify object that needs to be created"
// @Success 		200 		{object} 	dto.ResponseMessage
// @Failure 		400			{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 		404			{object} 	dto.ValidationErrorResponse "Email or Code does not exists."
// @Failure 		409			{object} 	dto.ValidationErrorResponse "Email is already verified."
// @Failure 		410			{object} 	dto.ValidationErrorResponse "Code has expired."
// @Failure 		500 		{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 			/auth/verify [post]
func (h *AuthHandler) Verify(c *fiber.Ctx) error {
	var req dto.VerifyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid body request.",
		})
	}

	if err := h.svc.Verify(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "VerifyEmail", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Your account has been successfully verified.",
	})
}

// ResendVerification	Resend email verification
// @Summary 			Resend email verification
// @Description 		Resends the verification code to the user's email if it hasn't been verified yet.
// @Tags        		auth
// @Accept 				json
// @Produce 			json
// @Param 				resend	 	body		dto.ResendVerificationRequest true "resend object that needs to be created"
// @Success 			200 		{object} 	dto.ResponseMessage
// @Failure 			400			{object} 	dto.ValidationErrorResponse "Invalid request"
// @Failure 			404			{object} 	dto.ValidationErrorResponse "Email does not exists."
// @Failure 			409			{object} 	dto.ValidationErrorResponse "Email is already verified."
// @Failure 			500 		{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 				/auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
	var req dto.ResendVerificationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseError{
			Message: "Invalid body request.",
		})
	}

	if err := h.svc.ResendVerification(c.Context(), req); err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "ResendVerification", err)
	}

	return c.JSON(dto.ResponseMessage{
		Message: "Verification email resent successfully. Please check your inbox.",
	})
}

// VerificationStatus	Check email verification status
// @Summary 			Check email verification status
// @Description 		Checks whether the user's email has been verified.
// @Tags        		auth
// @Accept 				json
// @Produce 			json
// @Param 				email 		query 		string true "User email"
// @Success 			200 		{object} 	dto.ResponseWithData[any]
// @Failure 			404			{object} 	dto.ValidationErrorResponse "Email does not exists."
// @Failure 			500 		{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 				/auth/verification-status [get]
func (h *AuthHandler) VerificationStatus(c *fiber.Ctx) error {
	email := c.Query("email")

	if _, err := mail.ParseAddress(email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ResponseError{
			Message: "Invalid email format.",
		})
	}

	status, err := h.svc.VerificationStatus(c.Context(), email)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "VerificationStatus", err)
	}

	return c.JSON(dto.ResponseWithData[any]{
		Data: fiber.Map{
			"email":    email,
			"verified": status,
		},
	})
}

// AuthMe				Get current authenticated user info
// @Summary				Get current authenticated user info
// @Description 		Returns profile information of the currently authenticated user based on the provided JWT token.
// @Tags        		auth
// @Accept 				json
// @Produce 			json
// @Success 			200 		{object} 	dto.ResponseWithData[dto.User]
// @Failure 			500 		{object} 	dto.InternalErrorResponse "Internal server error"
// @Router 				/auth/me [get]
func (h *AuthHandler) AuthMe(c *fiber.Ctx) error {
	userID := utils.GetUserId(c.Context())

	user, err := h.svc.AuthMe(c.Context(), userID)
	if err != nil {
		return errs.HandleHTTPError(c, h.log, "auth_handler", "Me", err)
	}

	return c.JSON(dto.ResponseWithData[dto.User]{
		Data: user,
	})
}
