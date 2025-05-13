package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type FavoriteRepository interface {
	FindFavoriteSongsByUserId(ctx context.Context, userId, pageSize, offset int) (songs []models.Song, err error)
	FindCountFavoriteSongsByUserId(ctx context.Context, userId int) (total int, err error)
	FindExistsFavoriteSongBySongId(ctx context.Context, userId, songId int) (exists bool, err error)
	StoreFavoriteSong(ctx context.Context, userId, songId int) (err error)
	DeleteFavoriteSong(ctx context.Context, userId, songId int) (err error)
}

type FavoriteService interface {
	GetAllFavoriteSongsByUserId(ctx context.Context, userId, pageSize, offset int) (songs []dto.Song, err error)
	GetCountFavoriteSongsByUserId(ctx context.Context, userId int) (total int, err error)
	CreateFavoriteSong(ctx context.Context, userId, songId int) (err error)
	DeleteFavoriteSong(ctx context.Context, userId, songId int) (err error)
}
