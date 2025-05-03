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
	authGroup.Post("/login", h.Auth.Login)
	authGroup.Post("/register", h.Auth.Register)
	authGroup.Put("/verify-email", h.Middleware.AuthRequired(), h.Auth.VerifyEmail)
	authGroup.Post("/resend-code", h.Middleware.AuthRequired(), h.Auth.ResendCodeEmailVerification)
	authGroup.Get("/profile", h.Middleware.AuthRequired(), h.Auth.Profile)

	v1.Get("/users", h.Middleware.AuthRequired(), h.User.GetUsers)
	v1.Get("/users/:id", h.Middleware.AuthRequired(), h.User.GetUser)
	v1.Put("/users/:id", h.Middleware.AuthRequired(), h.User.Update)
	v1.Delete("/users/:id", h.Middleware.AuthRequired(), h.User.Delete)

	return app
}
