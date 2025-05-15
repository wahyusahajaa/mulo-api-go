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
	// GetAll Return list of songs and total.
	//  Returns:
	//   200 OK: Success with lists and total.
	//   500 Internal Server Error: On Failure.
	GetAll(ctx context.Context, pageSize, offset int) (songs []dto.Song, total int, err error)

	// GetSongById Retrieve a song by Id.
	//  Returns:
	//   200 OK: on success with a song.
	//   400 Not Found if song does not exist.
	//   500 Internal Server Error: on failure.
	GetSongById(ctx context.Context, id int) (song dto.Song, err error)

	// CreateSong insert a new song.
	//  Returns:
	//   201 Created: on success.
	//   400 Bad Request: on validation failure.
	//   500 Internal Server Error: on failure.
	CreateSong(ctx context.Context, req dto.CreateSongRequest) (err error)

	// UpdateSong update an existing song by ID.
	//  Returns:
	//   200 OK: on success.
	//   400 Bad Request: on validation failure.
	//   404 Not Found: if song is missing.
	//   500 Internal Server Error: on failure.
	UpdateSong(ctx context.Context, req dto.CreateSongRequest, id int) (err error)

	// DeleteSong remove a song by ID.
	//  Returns:
	//   200 OK: on success.
	//   404 Not Found: if song is missing.
	//   500 Internal Server Error: on failure.
	DeleteSong(ctx context.Context, id int) (err error)
}
