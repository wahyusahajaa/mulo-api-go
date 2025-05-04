package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type ArtistRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (artists []models.Artist, err error)
	FindByArtistIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error)
	Count(ctx context.Context) (total int, err error)
	Store(ctx context.Context, name, slug string, image []byte) (err error)
	FindDuplicateArtistBySlug(ctx context.Context, slug string) (exists bool, err error)
	FindArtistById(ctx context.Context, artistId int) (artist *models.Artist, err error)
	Update(ctx context.Context, name, slug string, image []byte, artistId int) (err error)
	Delete(ctx context.Context, artistId int) (err error)
}

type ArtistService interface {
	GetAll(ctx context.Context, pageSize, offset int) (artists []models.Artist, err error)
	GetArtistByIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error)
	GetCount(ctx context.Context) (total int, err error)
	Create(ctx context.Context, name, slug string, image []byte) (err error)
	CheckDuplicateArtistBySlug(ctx context.Context, slug string) (exists bool, err error)
	GetArtistById(ctx context.Context, artistId int) (artist *models.Artist, err error)
	Update(ctx context.Context, name, slug string, image []byte, artistId int) (err error)
	Delete(ctx context.Context, artistId int) (err error)
}
