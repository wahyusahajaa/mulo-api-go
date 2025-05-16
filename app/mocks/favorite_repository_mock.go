package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type MockFavoriteRepository struct {
	mock.Mock
}

func (m *MockFavoriteRepository) DeleteFavoriteSong(ctx context.Context, userID int, songID int) (err error) {
	args := m.Called(ctx, userID, songID)

	return args.Error(0)
}

func (m *MockFavoriteRepository) FindCountFavoriteSongsByUserID(ctx context.Context, userID int) (total int, err error) {
	args := m.Called(ctx, userID)

	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *MockFavoriteRepository) FindExistsFavoriteSongBySongID(ctx context.Context, userID int, songID int) (exists bool, err error) {
	args := m.Called(ctx, userID, songID)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockFavoriteRepository) FindFavoriteSongsByUserID(ctx context.Context, userID int, pageSize int, offset int) (songs []models.Song, err error) {
	args := m.Called(ctx, userID, pageSize, offset)

	if args.Get(0) != nil {
		songs = args.Get(0).([]models.Song)
	}

	return songs, args.Error(1)
}

func (m *MockFavoriteRepository) StoreFavoriteSong(ctx context.Context, userID int, songID int) (err error) {
	args := m.Called(ctx, userID, songID)

	return args.Error(0)
}
