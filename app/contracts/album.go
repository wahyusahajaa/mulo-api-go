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
}

type AlbumService interface {
	GetAll(ctx context.Context, pageSize, offset int) (albums []dto.Album, err error)
	GetAlbumById(ctx context.Context, id int) (album dto.Album, err error)
	GetCount(ctx context.Context) (total int, err error)
	CreateAlbum(ctx context.Context, req dto.CreateAlbumRequest) (err error)
	UpdateAlbum(ctx context.Context, req dto.CreateAlbumRequest, id int) (err error)
	DeleteAlbum(ctx context.Context, id int) (err error)
}
