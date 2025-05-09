package contracts

import (
	"context"

	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type GenreRepository interface {
	FindAll(ctx context.Context, pageSize, offset int) (genres []models.Genre, err error)
	FindCount(ctx context.Context) (total int, err error)
	FindExistsGenreById(ctx context.Context, id int) (exists bool, err error)
	FindGenreById(ctx context.Context, id int) (genre *models.Genre, err error)
	Store(ctx context.Context, input models.CreateGenreInput) (err error)
	Update(ctx context.Context, input models.CreateGenreInput, id int) (err error)
	Delete(ctx context.Context, id int) (err error)
	StoreArtistGenre(ctx context.Context, artistId, genreId int) (err error)
	FindExistsArtistGenreByGenreId(ctx context.Context, artistId, genreId int) (exists bool, err error)
	FindArtistGenres(ctx context.Context, artistId, pageSize, offset int) (genres []models.Genre, err error)
	DeleteArtistGenre(ctx context.Context, artistId, genreId int) (err error)
	StoreSongGenre(ctx context.Context, songId, genreId int) (err error)
	FindExistsSongGenreByGenreId(ctx context.Context, songId, genreId int) (exists bool, err error)
	FindSongGenres(ctx context.Context, songId, pageSize, offset int) (genres []models.Genre, err error)
	DeleteSongGenre(ctx context.Context, songId, genreId int) (err error)
	FindAllArtists(ctx context.Context, genreId, pageSize, offset int) (artists []models.Artist, err error)
	FindCountArtists(ctx context.Context, genreId int) (total int, err error)
	FindAllSongs(ctx context.Context, genreId, pageSize, offset int) (songs []models.Song, err error)
	FindCountSongs(ctx context.Context, genreId int) (total int, err error)
}

type GenreService interface {
	GetAll(ctx context.Context, pageSize, offset int) (genres []dto.Genre, err error)
	GetCount(ctx context.Context) (total int, err error)
	GetGenreById(ctx context.Context, id int) (genre dto.Genre, err error)
	CreateGenre(ctx context.Context, req dto.CreateGenreRequest) (err error)
	UpdateGenre(ctx context.Context, req dto.CreateGenreRequest, id int) (err error)
	DeleteGenre(ctx context.Context, id int) (err error)
	CreateArtistGenre(ctx context.Context, artistId, genreId int) (err error)
	GetArtistGenres(ctx context.Context, artistId, pageSize, offset int) (genres []dto.Genre, err error)
	DeleteArtistGenre(ctx context.Context, artistId, genreId int) (err error)
	CreateSongGenre(ctx context.Context, songId, genreId int) (err error)
	GetSongGenres(ctx context.Context, songId, pageSize, offset int) (genres []dto.Genre, err error)
	DeleteSongGenre(ctx context.Context, songId, genreId int) (err error)
	GetAllArtists(ctx context.Context, genreId, pageSize, offset int) (artists []dto.Artist, err error)
	GetCountArtists(ctx context.Context, genreId int) (total int, err error)
	GetAllSongs(ctx context.Context, genreId, pageSize, offset int) (songs []dto.Song, err error)
	GetCountSongs(ctx context.Context, genreId int) (total int, err error)
}
