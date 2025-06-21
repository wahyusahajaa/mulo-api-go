package main

import (
	"log"

	"github.com/wahyusahajaa/mulo-api-go/app/di"
)

// @title Mulo Music Streaming API
// @version 1.0
// @description This documentation for access Mulo Music Streaming
// @contact.name The Developer
// @contact.email wahyusahaja.official@gmail.com
// @host api.mulo.craftedfolio.my.id
// @BasePath /v1
// @schemes https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type: Bearer token
func main() {
	// Initialized the application with all dependencies
	app, err := di.InitializedApp()
	if err != nil {
		log.Fatalf("failed to Initialized app: %v", err)
	}

	if err := app.App.Listen(":" + app.Config.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
