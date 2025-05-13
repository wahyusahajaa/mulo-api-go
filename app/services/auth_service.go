package services

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/jwt"
	"github.com/wahyusahajaa/mulo-api-go/pkg/resend"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
	"github.com/wahyusahajaa/mulo-api-go/pkg/verification"
)

type authService struct {
	repo            contracts.AuthRepository
	jwtSvc          jwt.JWTService
	verificationSvc verification.VerificationService
	resendSvc       resend.ResendService
	log             *logrus.Logger
}

func NewAuthService(repo contracts.AuthRepository, jwtSvc jwt.JWTService, verificationSvc verification.VerificationService, resendSvc resend.ResendService, log *logrus.Logger) contracts.AuthService {
	return &authService{
		repo:            repo,
		jwtSvc:          jwtSvc,
		verificationSvc: verificationSvc,
		resendSvc:       resendSvc,
		log:             log,
	}
}

func (svc *authService) Register(ctx context.Context, req dto.RegisterRequest) (err error) {
	// validation struct
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}

	// Check if email already exists
	exists, err := svc.repo.FindUserExistsByEmail(ctx, req.Email)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Register", err)
		return err
	}
	if exists {
		conflictErr := utils.ConflictError{Resource: "User", Field: "email", Value: req.Email}
		utils.LogWarn(svc.log, ctx, "auth_service", "Register", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	// Check if username already exists
	exists, err = svc.repo.FindUserExistsByUsername(ctx, req.Username)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Register", err)
		return err
	}
	if exists {
		conflictErr := utils.ConflictError{Resource: "User", Field: "username", Value: req.Username}
		utils.LogWarn(svc.log, ctx, "auth_service", "Register", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	// Generate verification code
	code, err := svc.verificationSvc.GenerateVerificationCode(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Register", err)
		return err
	}

	// transform req to input
	input := models.RegisterInput{
		Fullname: req.Fullname,
		Username: req.Username,
		Email:    req.Email,
		Password: utils.HashPassword(req.Password),
		Code:     code,
	}

	if err = svc.repo.Store(ctx, input); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Register", err)
		return err
	}

	// Send email code verification
	// go svc.resendService.SendEmailVerificationCode(req.Email, code)

	return
}

func (svc *authService) Login(ctx context.Context, req dto.LoginRequest) (token string, err error) {
	// validation struct
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return "", fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}

	// Get existing user by email
	user, err := svc.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Login", err)
		return "", err
	}
	if user == nil {
		notFoundErr := utils.NotFoundErrorWithCustomField{Resource: "User", Field: "email", Value: req.Email}
		utils.LogWarn(svc.log, ctx, "auth_service", "Login", notFoundErr)
		return "", fmt.Errorf("%w", notFoundErr)
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		notFoundErr := utils.NotFoundErrorWithCustomMsg{Message: "Invalid email or password."}
		utils.LogWarn(svc.log, ctx, "auth_service", "Login", notFoundErr)
		return "", fmt.Errorf("%w", notFoundErr)
	}

	// Generate token
	token, err = svc.jwtSvc.GenerateJWTToken(user.Id, user.Fullname, user.Username, user.Role, user.EmailVerifiedAt.Valid)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Login", err)
		return "", err
	}

	return
}

func (svc *authService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest, userId int) (err error) {
	// validation struct
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}

	// Check existing user verify by code
	userVerified, err := svc.repo.FindUserVerifiedByUserIdAndCode(ctx, userId, req.Code)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "VerifyEmail", err)
		return err
	}
	if userVerified == nil {
		notFoundErr := utils.NotFoundErrorWithCustomField{Resource: "Verification", Field: "code", Value: req.Code}
		utils.LogWarn(svc.log, ctx, "auth_service", "VerifyEmail", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check expired code
	if userVerified.ExpiredAt.Valid && userVerified.ExpiredAt.Time.Before(time.Now()) {
		goneErr := utils.GoneError{Resource: "Verification Code", Field: "code", Value: req.Code}
		utils.LogWarn(svc.log, ctx, "auth_service", "VerifyEmail", goneErr)
		return fmt.Errorf("%w", goneErr)
	}

	// Update user
	if err = svc.repo.UpdateUserVerifiedAt(ctx, userId); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "VerifyEmail", err)
		return err
	}

	return nil
}

func (svc *authService) ResendCode(ctx context.Context, userId int) (err error) {
	// Generate code
	code, err := svc.verificationSvc.GenerateVerificationCode(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "ResendCode", err)
		return err
	}

	if err := svc.repo.StoreUserVerifyCode(ctx, userId, code); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "ResendCode", err)
		return err
	}

	// TODO: SEND EMAIL CODE VERIFICATION IN HERE

	return
}
