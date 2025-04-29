package handlers

import "github.com/wahyusahajaa/mulo-api-go/app/middlewares"

type Handlers struct {
	Auth       *AuthHandler
	Middleware *middlewares.AuthMiddleware
}

func NewHandlers(auth *AuthHandler, middleware *middlewares.AuthMiddleware) *Handlers {
	return &Handlers{
		Auth:       auth,
		Middleware: middleware,
	}
}
