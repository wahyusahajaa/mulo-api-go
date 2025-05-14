package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/wahyusahajaa/mulo-api-go/app/di"
	_ "github.com/wahyusahajaa/mulo-api-go/docs"
)

// @title Mulo Music Streaming API
// @version 1.0
// @description This documentation for access Mulo Music Streaming
// @contact.name The Developer
// @contact.email wahyusahaja.official@gmail.com
// @host localhost:3000
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type: Bearer token
// @BasePath /v1
// @schemes http https
func main() {
	// Initialized the application with all dependencies
	app, err := di.InitializedApp()
	if err != nil {
		log.Fatalf("failed to Initialized app: %v", err)
	}

	app.App.Get("/swagger/*", swagger.HandlerDefault)

	// Start the server
	log.Printf("Server starting on port %s", app.Config.AppPort)
	if err := app.App.Listen(":" + app.Config.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// Ping godoc
// @Summary      Health Check
// @Description  Returns pong
// @Tags         health
// @Success      200 {object} string
// @Router       /ping [get]
func Ping(c *fiber.Ctx) error {
	return c.SendString("pong")
}
