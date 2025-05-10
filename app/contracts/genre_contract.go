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
	// GetAll returns a list of genres and total count.
	//  Returns:
	//   200 OK:: Success with list and total.
	//   500 Internal Server Error:: On failure.
	GetAll(ctx context.Context, pageSize, offset int) (genres []dto.Genre, total int, err error)

	// GetGenreById retrieves a genre by its ID.
	//  Returns:
	//   200 OK: on success.
	//   400 Not Found if genre does not exist.
	//   500 Internal Server Error: on failure.
	GetGenreById(ctx context.Context, id int) (genre dto.Genre, err error)

	// CreateGenre insert a new genre.
	//  Returns:
	//   201 Created: on success.
	//   400 Bad Request: on validation failure.
	//   500 Internal Server Error: on failure.
	CreateGenre(ctx context.Context, req dto.CreateGenreRequest) (err error)

	// UpdateGenre update an existing genre by ID.
	//  Returns:
	//   200 OK: on success.
	//   400 Bad Request: on validation failure.
	//   404 Not Found: if genre is missing.
	//   500 Internal Server Error: on failure.
	UpdateGenre(ctx context.Context, req dto.CreateGenreRequest, id int) (err error)

	// DeleteGenre remove a genre by ID.
	//  Returns:
	//   204 No Content on success.
	//   404 Not Found: if genre is missing.
	//   500 Internal Server Error: on failure.
	DeleteGenre(ctx context.Context, id int) (err error)

	// CreateArtistGenre insert a new artist genre.
	//  Returns:
	//   201 Created: on success.
	//   404 Not Found: if artist and genre is missing.
	//   409 Conflict: if genre already exist with the same artist.
	//   500 Internal Server Error: on failure.
	CreateArtistGenre(ctx context.Context, artistId, genreId int) (err error)

	// GetArtistGenres returns a lists of artist genres.
	//  Returns:
	//   200 OK: with the lists.
	//   500 Internal Server Error: on failure.
	GetArtistGenres(ctx context.Context, artistId, pageSize, offset int) (genres []dto.Genre, err error)

	// DeleteArtistGenre remove genre from artist genre by artist and genre ID.
	//  Returns:
	//   204 No Content: on success.
	//   404 Not Found: if artist genre is missing.
	//   500 Internal Server Error: on failure.
	DeleteArtistGenre(ctx context.Context, artistId, genreId int) (err error)

	// CreateSongGenre insert a new song genre.
	//  Returns:
	//   201 Created: on success.
	//   404 Not Found: if song and genre is missing.
	//   409 Conflict: if genre already exist with the same genre.
	//   500 Internal Server Error: on failure.
	CreateSongGenre(ctx context.Context, songId, genreId int) (err error)

	// GetSongGenres returns a lists song genres
	//  Returns:
	//   200 OK: with the lists.
	//   500 Internal Server Error: on failure.
	GetSongGenres(ctx context.Context, songId, pageSize, offset int) (genres []dto.Genre, err error)

	// DeleteSongGenre remove genre from song genre by song and genre ID.
	//  Returns:
	//   204 No Content: on success.
	//   404 Not Found: if song genre is missing.
	//   500 Internal Server Error: on failure.
	DeleteSongGenre(ctx context.Context, songId, genreId int) (err error)

	// GetAllArtists returns a list of artists and total count by genre ID.
	//  Returns:
	//   200 OK: with the lists and total.
	//   500 Internal Server Error: on failure.
	GetAllArtists(ctx context.Context, genreId, pageSize, offset int) (artists []dto.Artist, total int, err error)

	// GetAllSongs returns a list of songs and total count by genre ID.
	//  Returns:
	//   200 OK: with the lists and total.
	//   500 Internal Server Error: on failure.
	GetAllSongs(ctx context.Context, genreId, pageSize, offset int) (songs []dto.Song, total int, err error)
}
