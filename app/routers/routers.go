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

	// Albums endpoint
	v1Protected.Get("/albums", h.Album.GetAlbums)
	v1Protected.Get("/albums/:id", h.Album.GetAlbum)
	v1Protected.Post("/albums", h.Album.CreateAlbum)
	v1Protected.Put("/albums/:id", h.Album.UpdateAlbum)
	v1Protected.Delete("/albums/:id", h.Album.DeleteAlbum)

	return app
}
