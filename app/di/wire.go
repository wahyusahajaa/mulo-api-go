//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/wahyusahajaa/mulo-api-go/app/config"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/handlers"
	"github.com/wahyusahajaa/mulo-api-go/app/middlewares"
	"github.com/wahyusahajaa/mulo-api-go/app/repositories"
	"github.com/wahyusahajaa/mulo-api-go/app/routers"
	"github.com/wahyusahajaa/mulo-api-go/app/services"
	"github.com/wahyusahajaa/mulo-api-go/pkg/logger"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type AppContainer struct {
	App    *fiber.App
	Config *config.Config
}

var authSet = wire.NewSet(
	utils.NewJWTService,
	utils.NewResendService,
	repositories.NewAuthRepository,
	services.NewAuthService,
	utils.NewVerification,
	handlers.NewAuthHandler,
)

var userSet = wire.NewSet(
	repositories.NewUserRepository,
	services.NewUserService,
	handlers.NewUserHandler,
)

var artistSet = wire.NewSet(
	repositories.NewArtistRepository,
	services.NewArtistService,
	handlers.NewArtistHandler,
)

var albumSet = wire.NewSet(
	repositories.NewAlbumRepository,
	services.NewAlbumService,
	handlers.NewAlbumHandler,
)

func InitializedApp() (*AppContainer, error) {
	wire.Build(
		logger.NewLogger,
		middlewares.FiberLogger,
		config.NewConfig,
		database.NewDB,
		authSet,
		userSet,
		artistSet,
		albumSet,
		middlewares.NewAuthMiddleware,
		handlers.NewHandlers,
		routers.ProviderFiberApp,
		wire.Struct(new(AppContainer), "*"),
	)

	return nil, nil
}
