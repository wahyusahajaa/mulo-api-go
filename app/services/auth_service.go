package services

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/jwt"
	"github.com/wahyusahajaa/mulo-api-go/pkg/resend"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
	"github.com/wahyusahajaa/mulo-api-go/pkg/verification"
)

type authService struct {
	authRepo        contracts.AuthRepository
	userRepo        contracts.UserRepository
	jwtSvc          jwt.JWTService
	verificationSvc verification.VerificationService
	resendSvc       resend.ResendService
	log             *logrus.Logger
}

func NewAuthService(
	authRepo contracts.AuthRepository,
	userRepo contracts.UserRepository,
	jwtSvc jwt.JWTService,
	verificationSvc verification.VerificationService,
	resendSvc resend.ResendService,
	log *logrus.Logger,
) contracts.AuthService {
	return &authService{
		authRepo:        authRepo,
		userRepo:        userRepo,
		jwtSvc:          jwtSvc,
		verificationSvc: verificationSvc,
		resendSvc:       resendSvc,
		log:             log,
	}
}

func (svc *authService) Register(ctx context.Context, req dto.RegisterRequest) (err error) {
	// validation struct
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	// Check if email already exists
	exists, err := svc.userRepo.FindUserExistsByEmail(ctx, req.Email)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Register", err)
		return err
	}
	if exists {
		conflictErr := errs.NewConflictError("User", "email", req.Email)
		utils.LogWarn(svc.log, ctx, "auth_service", "Register", conflictErr)
		return conflictErr
	}

	// Check if username already exists
	exists, err = svc.userRepo.FindUserExistsByUsername(ctx, req.Username)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Register", err)
		return err
	}
	if exists {
		conflictErr := errs.NewConflictError("User", "username", req.Username)
		utils.LogWarn(svc.log, ctx, "auth_service", "Register", conflictErr)
		return conflictErr
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

	if err = svc.authRepo.Store(ctx, input); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Register", err)
		return err
	}

	// Send email code verification
	// go svc.resendSvc.SendEmailVerificationCode(req.Email, code)

	return
}

func (svc *authService) Login(ctx context.Context, req dto.LoginRequest) (accessToken, refreshToken string, err error) {
	// validation struct
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return "", "", errs.NewBadRequestError("validation failed", errorsMap)
	}

	// Get existing user by email
	user, err := svc.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Login", err)
		return "", "", err
	}
	if user == nil {
		notFoundErr := errs.NewNotFoundError("User", "email", req.Email)
		utils.LogWarn(svc.log, ctx, "auth_service", "Login", notFoundErr)
		return "", "", notFoundErr
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		notFoundErr := errs.NewNotFoundErrorWithMsg("Password mismatch. Try again.")
		utils.LogWarn(svc.log, ctx, "auth_service", "Login", notFoundErr)
		return "", "", notFoundErr
	}

	if !user.EmailVerifiedAt.Valid {
		forbiddenErr := errs.NewForbiddenError("Access denied. Please verify your email to continue.")
		utils.LogWarn(svc.log, ctx, "auth_service", "login", forbiddenErr)
		return "", "", forbiddenErr
	}

	accessToken, refreshToken, err = svc.jwtSvc.GenerateTokens(user.Id, user.Username, user.Role)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Login", err)
		return "", "", err
	}

	// Store Refresh Token
	refreshTokenExpires := time.Now().Add(7 * 24 * time.Hour)
	if err := svc.authRepo.StoreRefreshToken(ctx, user.Id, refreshToken, refreshTokenExpires); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Login", err)
		return "", "", err
	}

	return
}

func (svc *authService) Verify(ctx context.Context, req dto.VerifyRequest) (err error) {
	// validation struct
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	// Check existing user by email
	user, err := svc.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Verify", err)
		return
	}

	if user == nil {
		nfErr := errs.NewNotFoundError("User", "email", req.Email)
		utils.LogWarn(svc.log, ctx, "auth_service", "Verify", nfErr)
		return nfErr
	}

	if user.EmailVerifiedAt.Valid {
		conflictErr := errs.NewConflictErrorWithMsg("Email is already verified.")
		utils.LogWarn(svc.log, ctx, "auth_service", "Verify", conflictErr)
		return conflictErr
	}

	// Check existing user verify by code
	userVerified, err := svc.userRepo.FindUserVerifiedByUserIDAndCode(ctx, user.Id, req.Code)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Verify", err)
		return err
	}
	if userVerified == nil {
		notFoundErr := errs.NewNotFoundError("Verification", "code", req.Code)
		utils.LogWarn(svc.log, ctx, "auth_service", "Verify", notFoundErr)
		return notFoundErr
	}

	// Check expired code
	if userVerified.ExpiredAt.Valid && userVerified.ExpiredAt.Time.Before(time.Now()) {
		goneErr := errs.NewGoneError("Verification Code", "code", req.Code)
		utils.LogWarn(svc.log, ctx, "auth_service", "Verify", goneErr)
		return goneErr
	}

	// Update user
	if err = svc.authRepo.UpdateUserVerifiedAt(ctx, user.Id); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Verify", err)
		return err
	}

	return nil
}

func (svc *authService) ResendVerification(ctx context.Context, req dto.ResendVerificationRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	user, err := svc.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "ResendVerification", err)
		return err
	}
	if user == nil {
		nfErr := errs.NewNotFoundError("User", "email", req.Email)
		utils.LogWarn(svc.log, ctx, "auth_service", "ResendVerification", nfErr)
		return nfErr
	}

	if user.EmailVerifiedAt.Valid {
		conflictErr := errs.NewConflictErrorWithMsg("Email is already verified.")
		utils.LogWarn(svc.log, ctx, "auth_service", "ResendVerification", conflictErr)
		return conflictErr
	}

	// Generate new code
	code, err := svc.verificationSvc.GenerateVerificationCode(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "ResendVerification", err)
		return err
	}

	if err := svc.authRepo.StoreUserVerifyCode(ctx, user.Id, code); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "ResendVerification", err)
		return err
	}

	// Send email verification
	// go svc.resendSvc.SendEmailVerificationCode(req.Email, code)

	return
}

func (svc *authService) VerificationStatus(ctx context.Context, email string) (status bool, err error) {
	user, err := svc.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "VerificationStatus", err)
		return false, err
	}
	if user == nil {
		nfErr := errs.NewNotFoundError("User", "email", email)
		utils.LogWarn(svc.log, ctx, "auth_service", "VerificationStatus", nfErr)
		return false, nfErr
	}

	status = user.EmailVerifiedAt.Valid

	return
}

func (svc *authService) AuthMe(ctx context.Context, userID int) (user dto.User, err error) {
	result, err := svc.userRepo.FindUserByUserID(ctx, userID)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Profile", err)
		return dto.User{}, err
	}
	if result == nil {
		nfErr := errs.NewNotFoundError("User", "id", userID)
		utils.LogWarn(svc.log, ctx, "auth_service", "AuthMe", nfErr)
		return dto.User{}, nfErr
	}

	user.Id = result.Id
	user.Fullname = result.Fullname
	user.Email = result.Email
	user.Username = result.Username
	user.Image = utils.ParseImageToJSON(result.Image)

	return
}

func (svc *authService) Refresh(ctx context.Context, token string) (accessToken string, refreshToken string, err error) {
	claims, err := svc.jwtSvc.ParseRefreshToken(token)
	if err != nil {
		utils.LogWarn(svc.log, ctx, "auth_service", "Refresh", err)
		return "", "", err
	}

	// Check if refresh token is valid or not revoked
	valid, err := svc.authRepo.IsRefreshTokenValid(ctx, claims.ID, token)
	if err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Refresh", err)
		return "", "", err
	}
	if !valid {
		unauthErr := errs.NewUnauthorizedError("Invalid refresh token or revoked.")
		utils.LogWarn(svc.log, ctx, "auth_service", "refresh", unauthErr)
		return "", "", unauthErr
	}

	// Update revoke status to TRUE or is revoked
	if err = svc.authRepo.UpdateRefreshToken(ctx, claims.ID, token); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Refresh", err)
		return "", "", err
	}

	// Generate access and refresh token
	accessToken, refreshToken, err = svc.jwtSvc.GenerateTokens(claims.ID, claims.Username, claims.UserRole)
	if err != nil {
		utils.LogWarn(svc.log, ctx, "auth_service", "Refresh", err)
		return "", "", err
	}

	// Store refresh token for next request refresh token while acces token has expired
	refreshTokenExpires := time.Now().Add(7 * 24 * time.Hour)
	if err := svc.authRepo.StoreRefreshToken(ctx, claims.ID, refreshToken, refreshTokenExpires); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Refresh", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (svc *authService) Logout(ctx context.Context, token string) (err error) {
	if err = svc.authRepo.DeleteRefreshToken(ctx, token); err != nil {
		utils.LogError(svc.log, ctx, "auth_service", "Logout", err)
		return
	}

	return
}
