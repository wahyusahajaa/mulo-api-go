package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type FavoriteRepository interface {
	FindFavoriteSongsByUserID(ctx context.Context, userID, pageSize, offset int) (songs []models.Song, err error)
	FindCountFavoriteSongsByUserID(ctx context.Context, userID int) (total int, err error)
	FindExistsFavoriteSongBySongID(ctx context.Context, userID, songID int) (exists bool, err error)
	StoreFavoriteSong(ctx context.Context, userID, songID int) (err error)
	DeleteFavoriteSong(ctx context.Context, userID, songID int) (err error)
}

type FavoriteService interface {
	// Get list of favorite songs
	GetFavoriteSongsByUserID(ctx context.Context, userID, pageSize, offset int) (songs []dto.Song, total int, err error)

	// Add song to favorite
	AddFavoriteSong(ctx context.Context, userID, songID int) (err error)

	// Remove song from favorite
	RemoveFavoriteSong(ctx context.Context, userID, songID int) (err error)
}
