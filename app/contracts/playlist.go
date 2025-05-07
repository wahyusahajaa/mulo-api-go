package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type PlaylistRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (playlists []models.Playlist, err error)
	FindById(ctx context.Context, id int) (playlist *models.Playlist, err error)
	FindCount(ctx context.Context) (total int, err error)
	FindExistsPlaylistById(ctx context.Context, id int) (exists bool, err error)
	Store(ctx context.Context, input models.CreatePlaylistInput) (err error)
	Update(ctx context.Context, input models.CreatePlaylistInput, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
}

type PlaylistService interface {
	GetAll(ctx context.Context, pageSize, offset int) (playlists []dto.Playlist, err error)
	GetPlaylistById(ctx context.Context, id int) (playlist dto.Playlist, err error)
	GetCount(ctx context.Context) (total int, err error)
	CreatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest) (err error)
	UpdatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest, id int) (err error)
	DeletePlaylist(ctx context.Context, id int) (err error)
}
