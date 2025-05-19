package models

type OAuthAccount struct {
	ID             int
	UserID         int
	Provider       string
	ProviderUserID string
}

type OAuthAccountInput struct {
	ID       string
	Username string
	Fullname string
	Email    string
	Image    []byte
	Provider string // github, google, facebook
}
