package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type MockSongRepository struct {
	mock.Mock
}

func (m *MockSongRepository) Delete(ctx context.Context, id int) (err error) {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockSongRepository) FindAll(ctx context.Context, pageSize int, offset int) (songs []models.Song, err error) {
	args := m.Called(ctx, pageSize, offset)

	if args.Get(0) != nil {
		songs = args.Get(0).([]models.Song)
	}

	return songs, args.Error(1)
}

func (m *MockSongRepository) FindCount(ctx context.Context) (total int, err error) {
	args := m.Called(ctx)

	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *MockSongRepository) FindExistsSongById(ctx context.Context, id int) (exists bool, err error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockSongRepository) FindSongById(ctx context.Context, id int) (song *models.Song, err error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		song = args.Get(0).(*models.Song)
	}

	return song, args.Error(1)
}

func (m *MockSongRepository) Store(ctx context.Context, input models.CreateSongInput) (err error) {
	args := m.Called(ctx, input)

	return args.Error(0)
}

func (m *MockSongRepository) Update(ctx context.Context, input models.CreateSongInput, id int) (err error) {
	args := m.Called(ctx, input, id)

	return args.Error(0)
}
