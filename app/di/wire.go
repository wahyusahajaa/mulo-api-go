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
	"github.com/wahyusahajaa/mulo-api-go/pkg/jwt"
	"github.com/wahyusahajaa/mulo-api-go/pkg/logger"
	"github.com/wahyusahajaa/mulo-api-go/pkg/resend"
	"github.com/wahyusahajaa/mulo-api-go/pkg/verification"
)

type AppContainer struct {
	App    *fiber.App
	Config *config.Config
}

var commonSet = wire.NewSet(
	jwt.NewJWTService,
	resend.NewResendService,
	verification.NewVerificationService,
)

var authSet = wire.NewSet(
	repositories.NewAuthRepository,
	services.NewAuthService,
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

var songSet = wire.NewSet(
	repositories.NewSongRepository,
	services.NewSongService,
	handlers.NewSongHandler,
)

var genreSet = wire.NewSet(
	repositories.NewGenreRepository,
	services.NewGenreService,
	handlers.NewGenreHandler,
)

var playlistSet = wire.NewSet(
	repositories.NewPlaylistRepository,
	services.NewPlaylistService,
	handlers.NewPlaylistHandler,
)

var favoriteSet = wire.NewSet(
	repositories.NewFavoriteRepository,
	services.NewFavoriteService,
	handlers.NewFavoriteHandler,
)

func InitializedApp() (*AppContainer, error) {
	wire.Build(
		logger.NewLogger,
		middlewares.FiberLogger,
		config.NewConfig,
		database.NewDB,
		commonSet,
		authSet,
		userSet,
		artistSet,
		albumSet,
		songSet,
		genreSet,
		playlistSet,
		favoriteSet,
		middlewares.NewAuthMiddleware,
		handlers.NewHandlers,
		routers.ProviderFiberApp,
		wire.Struct(new(AppContainer), "*"),
	)

	return nil, nil
}
