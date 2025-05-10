package mocks

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type MockGenreRepository struct {
	mock.Mock
}

// Delete implements contracts.GenreRepository.
func (m *MockGenreRepository) Delete(ctx context.Context, id int) (err error) {
	panic("unimplemented")
}

// DeleteArtistGenre implements contracts.GenreRepository.
func (m *MockGenreRepository) DeleteArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	panic("unimplemented")
}

// DeleteSongGenre implements contracts.GenreRepository.
func (m *MockGenreRepository) DeleteSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	panic("unimplemented")
}

// FindAllArtists implements contracts.GenreRepository.
func (m *MockGenreRepository) FindAllArtists(ctx context.Context, genreId int, pageSize int, offset int) (artists []models.Artist, err error) {
	panic("unimplemented")
}

// FindAllSongs implements contracts.GenreRepository.
func (m *MockGenreRepository) FindAllSongs(ctx context.Context, genreId int, pageSize int, offset int) (songs []models.Song, err error) {
	panic("unimplemented")
}

// FindArtistGenres implements contracts.GenreRepository.
func (m *MockGenreRepository) FindArtistGenres(ctx context.Context, artistId int, pageSize int, offset int) (genres []models.Genre, err error) {
	panic("unimplemented")
}

// FindCountArtists implements contracts.GenreRepository.
func (m *MockGenreRepository) FindCountArtists(ctx context.Context, genreId int) (total int, err error) {
	panic("unimplemented")
}

// FindCountSongs implements contracts.GenreRepository.
func (m *MockGenreRepository) FindCountSongs(ctx context.Context, genreId int) (total int, err error) {
	panic("unimplemented")
}

// FindExistsArtistGenreByGenreId implements contracts.GenreRepository.
func (m *MockGenreRepository) FindExistsArtistGenreByGenreId(ctx context.Context, artistId int, genreId int) (exists bool, err error) {
	panic("unimplemented")
}

// FindExistsGenreById implements contracts.GenreRepository.
func (m *MockGenreRepository) FindExistsGenreById(ctx context.Context, id int) (exists bool, err error) {
	panic("unimplemented")
}

// FindExistsSongGenreByGenreId implements contracts.GenreRepository.
func (m *MockGenreRepository) FindExistsSongGenreByGenreId(ctx context.Context, songId int, genreId int) (exists bool, err error) {
	panic("unimplemented")
}

// FindGenreById implements contracts.GenreRepository.
func (m *MockGenreRepository) FindGenreById(ctx context.Context, id int) (genre *models.Genre, err error) {
	panic("unimplemented")
}

// FindSongGenres implements contracts.GenreRepository.
func (m *MockGenreRepository) FindSongGenres(ctx context.Context, songId int, pageSize int, offset int) (genres []models.Genre, err error) {
	panic("unimplemented")
}

// Store implements contracts.GenreRepository.
func (m *MockGenreRepository) Store(ctx context.Context, input models.CreateGenreInput) (err error) {
	panic("unimplemented")
}

// StoreArtistGenre implements contracts.GenreRepository.
func (m *MockGenreRepository) StoreArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	panic("unimplemented")
}

// StoreSongGenre implements contracts.GenreRepository.
func (m *MockGenreRepository) StoreSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	panic("unimplemented")
}

// Update implements contracts.GenreRepository.
func (m *MockGenreRepository) Update(ctx context.Context, input models.CreateGenreInput, id int) (err error) {
	panic("unimplemented")
}

func (m *MockGenreRepository) FindAll(ctx context.Context, pageSize, offset int) (genres []models.Genre, err error) {
	args := m.Called(ctx, pageSize, offset)

	if args.Get(0) != nil {
		genres = args.Get(0).([]models.Genre)
	}

	err = args.Error(1)

	return
}

func (m *MockGenreRepository) FindCount(ctx context.Context) (total int, err error) {
	args := m.Called(ctx)

	count, ok := args.Get(0).(int)
	if !ok && args.Get(0) != nil {
		return 0, fmt.Errorf("invalid type for count")
	}

	return count, args.Error(1)
}
