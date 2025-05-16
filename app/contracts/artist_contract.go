package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type ArtistRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (artists []models.Artist, err error)
	FindByArtistIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error)
	FindExistsArtistBySlug(ctx context.Context, slug string) (exists bool, err error)
	FindExistsArtistById(ctx context.Context, id int) (exists bool, err error)
	FindArtistById(ctx context.Context, artistId int) (artist *models.Artist, err error)
	FindCount(ctx context.Context) (total int, err error)
	Store(ctx context.Context, input models.CreateArtistInput) (err error)
	Update(ctx context.Context, input models.CreateArtistInput, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
}

type ArtistService interface {
	GetAll(ctx context.Context, pageSize, offset int) (artists []dto.Artist, total int, err error)
	CreateArtist(ctx context.Context, req dto.CreateArtistRequest) (err error)
	GetArtistById(ctx context.Context, artistId int) (artist dto.Artist, err error)
	UpdateArtist(ctx context.Context, req dto.CreateArtistRequest, id int) (err error)
	DeleteArtist(ctx context.Context, id int) (err error)
}
