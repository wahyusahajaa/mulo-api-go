package oauth

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
)

type OAuthService interface {
	// GithubAccessToken for get access token from github by github code from redirect url
	GithubAccessToken(ctx context.Context, code string) (accessToken string, err error)
	// GithubUserInfo for get user info by access token
	GithubUserInfo(ctx context.Context, token string) (user *dto.GithubUser, err error)
	// GithubUserEmail for get user email, it's call while email on user info is missing or empty
	GithubUserEmail(ctx context.Context, token string) (email string, err error)
}
