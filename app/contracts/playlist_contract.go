package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type PlaylistRepository interface {
	FindAll(ctx context.Context, userRole string, userId, pageSize, offset int) (playlists []models.Playlist, err error)
	FindById(ctx context.Context, userRole string, userId, id int) (playlist *models.Playlist, err error)
	FindCount(ctx context.Context, userRole string, userId int) (total int, err error)
	FindExistsPlaylistById(ctx context.Context, userRole string, userId, id int) (exists bool, err error)
	Store(ctx context.Context, input models.CreatePlaylistInput) (err error)
	Update(ctx context.Context, input models.CreatePlaylistInput, id int) (err error)
	Delete(ctx context.Context, userRole string, userId, playlistId int) (err error)
	FindPlaylistSongs(ctx context.Context, playlistId, pageSize, offset int) (songs []models.Song, err error)
	StorePlaylistSong(ctx context.Context, playlistId, songId int) (err error)
	FindExistsPlaylistSong(ctx context.Context, playlistId, songId int) (exists bool, err error)
	DeletePlaylistSong(ctx context.Context, playlistId, songId int) (err error)
}

type PlaylistService interface {
	GetAll(ctx context.Context, userRole string, userId, pageSize, offset int) (playlists []dto.Playlist, total int, err error)
	GetPlaylistById(ctx context.Context, userRole string, userId, id int) (playlist dto.Playlist, err error)
	CreatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest) (err error)
	UpdatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest, userRole string, userId, id int) (err error)
	DeletePlaylist(ctx context.Context, userRole string, userId, id int) (err error)
	GetPlaylistSongs(ctx context.Context, userRole string, userId, playlistId, pageSize, offset int) (songs []dto.Song, err error)
	CreatePlaylistSong(ctx context.Context, userRole string, userId, playlistId, songId int) (err error)
	DeletePlaylistSong(ctx context.Context, userRole string, userId, playlistId, songId int) (err error)
}
