package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type MockPlaylistRepository struct {
	mock.Mock
}

func (m *MockPlaylistRepository) Delete(ctx context.Context, userRole string, userId int, playlistId int) (err error) {
	args := m.Called(ctx, userRole, userId, playlistId)

	return args.Error(0)
}

func (m *MockPlaylistRepository) DeletePlaylistSong(ctx context.Context, playlistId int, songId int) (err error) {
	args := m.Called(ctx, playlistId, songId)

	return args.Error(0)
}

func (m *MockPlaylistRepository) FindAll(ctx context.Context, userRole string, userId int, pageSize int, offset int) (playlists []models.Playlist, err error) {
	args := m.Called(ctx, userRole, userId, pageSize, offset)

	if args.Get(0) != nil {
		playlists = args.Get(0).([]models.Playlist)
	}

	return playlists, args.Error(1)
}

func (m *MockPlaylistRepository) FindById(ctx context.Context, userRole string, userId int, id int) (playlist *models.Playlist, err error) {
	args := m.Called(ctx, userRole, userId, id)

	if args.Get(0) != nil {
		playlist = args.Get(0).(*models.Playlist)
	}

	return playlist, args.Error(1)
}

func (m *MockPlaylistRepository) FindCount(ctx context.Context, userRole string, userId int) (total int, err error) {
	args := m.Called(ctx, userRole, userId)

	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *MockPlaylistRepository) FindExistsPlaylistById(ctx context.Context, userRole string, userId int, id int) (exists bool, err error) {
	args := m.Called(ctx, userRole, userId, id)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockPlaylistRepository) FindExistsPlaylistSong(ctx context.Context, playlistId int, songId int) (exists bool, err error) {
	args := m.Called(ctx, playlistId, songId)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockPlaylistRepository) FindPlaylistSongs(ctx context.Context, playlistId int, pageSize int, offset int) (songs []models.Song, err error) {
	args := m.Called(ctx, playlistId, pageSize, offset)

	if args.Get(0) != nil {
		songs = args.Get(0).([]models.Song)
	}

	return songs, args.Error(1)
}

func (m *MockPlaylistRepository) Store(ctx context.Context, input models.CreatePlaylistInput) (err error) {
	args := m.Called(ctx, input)

	return args.Error(0)
}

func (m *MockPlaylistRepository) StorePlaylistSong(ctx context.Context, playlistId int, songId int) (err error) {
	args := m.Called(ctx, playlistId, songId)

	return args.Error(0)
}

func (m *MockPlaylistRepository) Update(ctx context.Context, input models.CreatePlaylistInput, id int) (err error) {
	args := m.Called(ctx, input, id)

	return args.Error(0)
}
