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

func (m *MockGenreRepository) Delete(ctx context.Context, id int) (err error) {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGenreRepository) DeleteArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	args := m.Called(ctx, artistId, genreId)

	return args.Error(0)
}

func (m *MockGenreRepository) DeleteSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	args := m.Called(ctx, songId, genreId)

	return args.Error(0)
}

func (m *MockGenreRepository) FindAllArtists(ctx context.Context, genreId int, pageSize int, offset int) (artists []models.Artist, err error) {
	args := m.Called(ctx, genreId, pageSize, offset)

	if args.Get(0) != nil {
		artists = args.Get(0).([]models.Artist)
	}

	return artists, args.Error(1)
}

func (m *MockGenreRepository) FindAllSongs(ctx context.Context, genreId int, pageSize int, offset int) (songs []models.Song, err error) {
	args := m.Called(ctx, genreId, pageSize, offset)

	if args.Get(0) != nil {
		songs = args.Get(0).([]models.Song)
	}

	return songs, args.Error(1)
}

func (m *MockGenreRepository) FindArtistGenres(ctx context.Context, artistId int, pageSize int, offset int) (genres []models.Genre, err error) {
	args := m.Called(ctx, artistId, pageSize, offset)

	if args.Get(0) != nil {
		genres = args.Get(0).([]models.Genre)
	}

	return genres, args.Error(1)
}

func (m *MockGenreRepository) FindCountArtists(ctx context.Context, genreId int) (total int, err error) {
	args := m.Called(ctx, genreId)

	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *MockGenreRepository) FindCountSongs(ctx context.Context, genreId int) (total int, err error) {
	args := m.Called(ctx, genreId)

	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *MockGenreRepository) FindExistsArtistGenreByGenreId(ctx context.Context, artistId int, genreId int) (exists bool, err error) {
	args := m.Called(ctx, artistId, genreId)

	return args.Get(0).(bool), args.Error(1)
}

func (m *MockGenreRepository) FindExistsGenreById(ctx context.Context, id int) (exists bool, err error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockGenreRepository) FindExistsSongGenreByGenreId(ctx context.Context, songId int, genreId int) (exists bool, err error) {
	args := m.Called(ctx, songId, genreId)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockGenreRepository) FindGenreById(ctx context.Context, id int) (genre *models.Genre, err error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		genre = args.Get(0).(*models.Genre)
	}

	return genre, args.Error(1)
}

func (m *MockGenreRepository) FindSongGenres(ctx context.Context, songId int, pageSize int, offset int) (genres []models.Genre, err error) {
	args := m.Called(ctx, songId, pageSize, offset)

	if args.Get(0) != nil {
		genres = args.Get(0).([]models.Genre)
	}

	return genres, args.Error(1)
}

func (m *MockGenreRepository) Store(ctx context.Context, input models.CreateGenreInput) (err error) {
	args := m.Called(ctx, input)

	return args.Error(0)
}

func (m *MockGenreRepository) StoreArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	args := m.Called(ctx, artistId, genreId)

	return args.Error(0)
}

func (m *MockGenreRepository) StoreSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	args := m.Called(ctx, songId, genreId)

	return args.Error(0)
}

func (m *MockGenreRepository) Update(ctx context.Context, input models.CreateGenreInput, id int) (err error) {
	args := m.Called(ctx, input, id)

	return args.Error(0)
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
