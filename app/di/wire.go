//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
	"github.com/wahyusahajaa/mulo-api-go/app/handler"
	"github.com/wahyusahajaa/mulo-api-go/app/router"
)

type AppContainer struct {
	App    *fiber.App
	Config *config.Config
}

func InitializedApp() (*AppContainer, error) {
	wire.Build(
		config.NewConfig,
		handler.NewHandlers,
		router.ProviderFiberApp,
		wire.Struct(new(AppContainer), "*"),
	)

	return nil, nil
}
