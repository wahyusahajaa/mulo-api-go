package main

import (
	"log"

	"github.com/wahyusahajaa/mulo-api-go/app/di"
)

func main() {
	// Initialized the application with all dependencies
	app, err := di.InitializedApp()
	if err != nil {
		log.Fatalf("failed to Initialized app: %v", err)
	}

	// Start the server
	log.Printf("Server starting on port %s", app.Config.AppPort)
	if err := app.App.Listen(":" + app.Config.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
