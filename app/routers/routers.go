package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/handlers"
)

func ProviderFiberApp(h *handlers.Handlers, fiberLogger fiber.Handler) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Mulo Music Streaming",
	})

	app.Use(fiberLogger)

	v1 := app.Group("/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Mulo")
	})

	authGroup := v1.Group("auth")

	// Public routes
	authGroup.Post("/login", h.Auth.Login)
	authGroup.Post("/register", h.Auth.Register)

	authGroupProtected := authGroup.Use(h.Middleware.AuthRequired())
	v1Protected := v1.Use(h.Middleware.AuthRequired())

	authGroupProtected.Put("/verify-email", h.Auth.VerifyEmail)
	authGroupProtected.Post("/resend-code", h.Auth.ResendCodeEmailVerification)
	authGroupProtected.Get("/profile", h.Auth.Profile)

	// Users endpoint
	v1Protected.Get("/users", h.User.GetUsers)
	v1Protected.Get("/users/:id", h.User.GetUser)
	v1Protected.Put("/users/:id", h.User.Update)
	v1Protected.Delete("/users/:id", h.User.Delete)

	// Artists endpoint
	v1Protected.Get("/artists", h.Artist.GetArtists)
	v1Protected.Get("/artists/:id", h.Artist.GetArtist)
	v1Protected.Post("/artists", h.Artist.CreateArtist)
	v1Protected.Put("/artists/:id", h.Artist.UpdateArtist)
	v1Protected.Delete("/artists/:id", h.Artist.DeleteArtist)
	// Artists genres endpoint
	v1Protected.Get("/artists/:id/genres", h.Genre.GetArtistGenres)
	v1Protected.Post("/artists/:id/genres/:genreId", h.Genre.CreateArtistGenre)
	v1Protected.Delete("/artists/:id/genres/:genreId", h.Genre.DeleteArtistGenre)
	// Artist albums endpoint
	v1Protected.Get("/artists/:id/albums", h.Album.GetAlbumsByArtistId)

	// Albums endpoint
	v1Protected.Get("/albums", h.Album.GetAlbums)
	v1Protected.Get("/albums/:id", h.Album.GetAlbum)
	v1Protected.Post("/albums", h.Album.CreateAlbum)
	v1Protected.Put("/albums/:id", h.Album.UpdateAlbum)
	v1Protected.Delete("/albums/:id", h.Album.DeleteAlbum)

	// Songs Endpoint
	v1Protected.Get("/songs", h.Song.GetSongs)
	v1Protected.Get("/songs/:id", h.Song.GetSong)
	v1Protected.Post("/songs", h.Song.CreateSong)
	v1Protected.Put("/songs/:id", h.Song.UpdateSong)
	v1Protected.Delete("/songs/:id", h.Song.DeleteSong)
	// Songs genres endpoint
	v1Protected.Get("/songs/:id/genres", h.Genre.GetSongGenres)
	v1Protected.Post("/songs/:id/genres/:genreId", h.Genre.CreateSongGenre)
	v1Protected.Delete("/songs/:id/genres/:genreId", h.Genre.DeleteSongGenre)

	// Genres Endpoint
	v1Protected.Get("/genres", h.Genre.GetGenres)
	v1Protected.Get("/genres/:id", h.Genre.GetGenre)
	v1Protected.Post("/genres", h.Genre.CreateGenre)
	v1Protected.Put("/genres/:id", h.Genre.UpdateGenre)
	v1Protected.Delete("/genres/:id", h.Genre.DeleteGenre)

	// Playlists Endpoint
	v1Protected.Get("/playlists", h.Playlist.GetPlaylists)
	v1Protected.Get("/playlists/:id", h.Playlist.GetPlaylist)
	v1Protected.Post("/playlists", h.Playlist.CreatePlaylist)
	v1Protected.Put("/playlists/:id", h.Playlist.UpdatePlaylist)
	v1Protected.Delete("/playlists/:id", h.Playlist.DeletePlaylist)

	// Song Favorites Endpoints
	v1Protected.Get("/favorites/songs", h.Favorite.GetFavoriteSongsByUserId)
	v1Protected.Post("/favorites/songs/:songId", h.Favorite.CreateFavoriteSong)
	v1Protected.Delete("/favorites/songs/:songId", h.Favorite.DeleteFavoriteSong)

	return app
}
