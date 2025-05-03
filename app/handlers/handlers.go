package handlers

import "github.com/wahyusahajaa/mulo-api-go/app/middlewares"

type Handlers struct {
	Auth       *AuthHandler
	Middleware *middlewares.AuthMiddleware
	User       *UserHandler
}

func NewHandlers(
	auth *AuthHandler,
	middleware *middlewares.AuthMiddleware,
	user *UserHandler) *Handlers {
	return &Handlers{
		Auth:       auth,
		Middleware: middleware,
		User:       user,
	}
}
