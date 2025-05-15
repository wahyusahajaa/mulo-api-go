package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type AlbumRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (albums []models.Album, err error)
	FindAlbumById(ctx context.Context, id int) (album *models.AlbumWithArtist, err error)
	FindCount(ctx context.Context) (total int, err error)
	FindExistsAlbumById(ctx context.Context, id int) (exists bool, err error)
	FindExistsAlbumBySlug(ctx context.Context, slug string) (exists bool, err error)
	Store(ctx context.Context, input models.CreateAlbumInput) (err error)
	Update(ctx context.Context, input models.CreateAlbumInput, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
	FindAlbumsByArtistId(ctx context.Context, artistId int) (albums []models.Album, err error)
}

type AlbumService interface {
	// GetAll Return list of albums and total.
	//  Returns:
	//   200 OK: Success with lists and total.
	//   500 Internal Server Error: On Failure.
	GetAll(ctx context.Context, pageSize, offset int) (albums []dto.AlbumWithArtist, total int, err error)

	// GetAlbumById Retrieve a album by Id.
	//  Returns:
	//   200 OK: on success with a album.
	//   400 Not Found: album does not exists.
	//   500 Internal Server Error: on failure.
	GetAlbumById(ctx context.Context, id int) (album dto.AlbumWithArtist, err error)

	// CreateAlbum insert a new album.
	//  Returns:
	//   201 Created: on success.
	//   400 Bad Request: on validation failure.
	//   404 Not Found: artist does not exists.
	//   409 Conflict: album name or slug already exists.
	//   500 Internal Server Error: on failure.
	CreateAlbum(ctx context.Context, req dto.CreateAlbumRequest) (err error)

	// UpdateAlbum update an existing album by ID.
	//  Returns:
	//   200 OK: on success.
	//   400 Bad Request: on validation failure.
	//   404 Not Found: album or artist does not exists.
	//   409 Conflict: album name or slug already exists.
	//   500 Internal Server Error: on failure.
	UpdateAlbum(ctx context.Context, req dto.CreateAlbumRequest, id int) (err error)

	// DeleteAlbum remove a album by ID.
	//  Returns:
	//   200 OK: on success.
	//   404 Not Found: album does not exists.
	//   500 Internal Server Error: on failure.
	DeleteAlbum(ctx context.Context, id int) (err error)

	// GetAll Return list of albums by artist id.
	//  Returns:
	//   200 OK: Success with lists.
	//   500 Internal Server Error: On Failure.
	GetAlbumsByArtistId(ctx context.Context, artistId int) (albums []dto.Album, err error)
}
