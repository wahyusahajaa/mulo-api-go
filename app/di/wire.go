//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
	"github.com/wahyusahajaa/mulo-api-go/app/handlers"
	"github.com/wahyusahajaa/mulo-api-go/app/routers"
)

type AppContainer struct {
	App    *fiber.App
	Config *config.Config
}

func InitializedApp() (*AppContainer, error) {
	wire.Build(
		config.NewConfig,
		handlers.NewHandlers,
		routers.ProviderFiberApp,
		wire.Struct(new(AppContainer), "*"),
	)

	return nil, nil
}
