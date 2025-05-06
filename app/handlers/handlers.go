package handlers

import "github.com/wahyusahajaa/mulo-api-go/app/middlewares"

type Handlers struct {
	Auth       *AuthHandler
	Middleware *middlewares.AuthMiddleware
	User       *UserHandler
	Artist     *ArtistHandler
	Album      *AlbumHandler
	Song       *SongHandler
	Genre      *GenreHandler
}

func NewHandlers(
	auth *AuthHandler,
	middleware *middlewares.AuthMiddleware,
	user *UserHandler,
	artist *ArtistHandler,
	album *AlbumHandler,
	song *SongHandler,
	genre *GenreHandler,
) *Handlers {
	return &Handlers{
		Auth:       auth,
		Middleware: middleware,
		User:       user,
		Artist:     artist,
		Album:      album,
		Song:       song,
		Genre:      genre,
	}
}
