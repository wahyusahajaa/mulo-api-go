package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type oauthService struct {
	githubClientID     string
	githubClientSecret string
	log                *logrus.Logger
}

func NewOauthService(conf *config.Config, log *logrus.Logger) OAuthService {
	return &oauthService{
		githubClientID:     conf.GithubClientID,
		githubClientSecret: conf.GithubClientSecret,
		log:                log,
	}
}

func (svc *oauthService) GithubAccessToken(ctx context.Context, code string) (accessToken string, err error) {
	var client = &http.Client{Timeout: 10 * time.Second}
	var url = `https://github.com/login/oauth/access_token`

	reqBody := fmt.Sprintf("client_id=%s&client_secret=%s&code=%s", svc.githubClientID, svc.githubClientSecret, code)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "application/json")
	req.URL.RawQuery = reqBody

	if err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubAccessToken", err)
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubAccessToken", err)
		return "", err
	}
	defer res.Body.Close()

	var resData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(res.Body).Decode(&resData); err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubAccessToken", err)
		return "", err
	}

	return resData.AccessToken, nil
}

func (svc *oauthService) GithubUserInfo(ctx context.Context, token string) (githubUser *dto.GithubUser, err error) {
	var client = &http.Client{Timeout: 10 * time.Second}
	var url = `https://api.github.com/user`

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")

	if err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubUserInfo", err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubUserInfo", err)
		return nil, err
	}
	defer res.Body.Close()

	githubUser = &dto.GithubUser{}
	if err := json.NewDecoder(res.Body).Decode(&githubUser); err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubUserInfo", err)
		return nil, err
	}

	return githubUser, nil
}

func (svc *oauthService) GithubUserEmail(ctx context.Context, token string) (email string, err error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://api.github.com/user/emails"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/json")
	if err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubUserEmail", err)
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubUserEmail", err)
		return "", err
	}

	var emails []dto.GithubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		utils.LogError(svc.log, ctx, "oauthGithub_service", "GithubUserEmail", err)
		return "", err
	}

	return emails[0].Email, nil
}
