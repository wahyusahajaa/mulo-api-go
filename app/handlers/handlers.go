package handlers

import "github.com/wahyusahajaa/mulo-api-go/app/middlewares"

type Handlers struct {
	Auth       *AuthHandler
	Middleware *middlewares.AuthMiddleware
	User       *UserHandler
	Artist     *ArtistHandler
	Album      *AlbumHandler
}

func NewHandlers(
	auth *AuthHandler,
	middleware *middlewares.AuthMiddleware,
	user *UserHandler,
	artist *ArtistHandler,
	album *AlbumHandler,
) *Handlers {
	return &Handlers{
		Auth:       auth,
		Middleware: middleware,
		User:       user,
		Artist:     artist,
		Album:      album,
	}
}
