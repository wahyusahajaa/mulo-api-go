package contracts

import (
	"context"
	"time"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type AuthRepository interface {
	Store(ctx context.Context, input models.RegisterInput) (err error)
	StoreUserVerifyCode(ctx context.Context, userId int, code string) (err error)
	UpdateUserVerifiedAt(ctx context.Context, userId int) (err error)

	IsRefreshTokenValid(ctx context.Context, userID int, token string) (exists bool, err error)
	StoreRefreshToken(ctx context.Context, userID int, token string, expiredAt time.Time) (err error)
	UpdateRefreshToken(ctx context.Context, userID int, token string) (err error)
	DeleteRefreshToken(ctx context.Context, token string) (err error)
	StoreUserWithOAuthAccount(ctx context.Context, input models.OAuthAccountInput) (userID int, err error)
	StoreOAuthAccount(ctx context.Context, userID int, providerID, providerUserID string) (err error)
	FindOAuthAccount(ctx context.Context, provider, providerUserID string) (*models.OAuthAccount, error)
	FindExistsOauthAccount(ctx context.Context, userID int) (exists bool, err error)
}

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (err error)
	Login(ctx context.Context, req dto.LoginRequest) (accessToken, refreshToken string, err error)
	Verify(ctx context.Context, req dto.VerifyRequest) (err error)
	ResendVerification(ctx context.Context, req dto.ResendVerificationRequest) (err error)
	VerificationStatus(ctx context.Context, email string) (status bool, err error)
	AuthMe(ctx context.Context, userID int) (user dto.User, err error)
	Refresh(ctx context.Context, token string) (accessToken, refreshToken string, err error)
	Logout(ctx context.Context, token string) (err error)

	// OAuthGithubCallback: Login or Register with github oAuth2
	//  Flows:
	//   Check user by email:
	//    if user already exist -> check oauth_account by user_id
	//     if oauth_account does not exists -> create oauth_accounts -> generate tokens
	//     if oauth_account already exists -> generate tokens
	//    if user does not exist -> create user and oauth_accounts -> generate tokens
	OAuthGithubCallback(ctx context.Context, req dto.GithubReq) (accessToken, refreshToken string, err error)
}
