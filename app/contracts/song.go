package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type SongRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (songs []models.Song, err error)
	FindCount(ctx context.Context) (total int, err error)
	FindSongById(ctx context.Context, id int) (song *models.Song, err error)
	FindExistsSongById(ctx context.Context, id int) (exists bool, err error)
	Store(ctx context.Context, input models.CreateSongInput) (err error)
	Update(ctx context.Context, input models.CreateSongInput, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
}

type SongService interface {
	GetAll(ctx context.Context, pageSize, offset int) (songs []dto.Song, err error)
	GetCount(ctx context.Context) (total int, err error)
	GetSongById(ctx context.Context, id int) (song dto.Song, err error)
	CreateSong(ctx context.Context, req dto.CreateSongRequest) (err error)
	UpdateSong(ctx context.Context, req dto.CreateSongRequest, id int) (err error)
	DeleteSong(ctx context.Context, id int) (err error)
}
