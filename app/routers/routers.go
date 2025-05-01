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
	authGroup.Post("/verify-email", h.Middleware.AuthRequired(), h.Auth.VerifyEmail)
	authGroup.Post("/resend-code", h.Middleware.AuthRequired(), h.Auth.ResendCodeEmailVerification)
	authGroup.Get("/me", h.Middleware.AuthRequired(), h.Auth.Me)

	return app
}
