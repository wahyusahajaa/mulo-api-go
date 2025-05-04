package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type AlbumRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (albums []models.Album, err error)
	FindAlbumById(ctx context.Context, id int) (album *models.AlbumWithArtist, err error)
	Count(ctx context.Context) (total int, err error)
	Store(ctx context.Context, artistId int, name, slug string, image []byte) (err error)
	FindDuplicateAlbumBySlug(ctx context.Context, slug string) (exists bool, err error)
	Update(ctx context.Context, artistId int, name, slug string, image []byte, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
}

type AlbumService interface {
	GetAll(ctx context.Context, pageSize, offset int) (albums []models.Album, err error)
	GetAlbumById(ctx context.Context, id int) (album *models.AlbumWithArtist, err error)
	GetCount(ctx context.Context) (total int, err error)
	Create(ctx context.Context, artistId int, name, slug string, image []byte) (err error)
	CheckDuplicateAlbumBySlug(ctx context.Context, slug string) (exists bool, err error)
	UpdateAlbum(ctx context.Context, artistId int, name, slug string, image []byte, id int) (err error)
	DeleteAlbum(ctx context.Context, id int) (err error)
}
