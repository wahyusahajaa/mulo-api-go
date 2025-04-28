package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wahyusahajaa/mulo-api-go/app/handlers"
)

func ProviderFiberApp(h *handlers.Handlers) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Mulo Music Streaming",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Mulo")
	})

	return app
}
