package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type MockAlbumRepository struct {
	mock.Mock
}

func (m *MockAlbumRepository) Delete(ctx context.Context, id int) (err error) {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockAlbumRepository) FindAlbumById(ctx context.Context, id int) (album *models.AlbumWithArtist, err error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		album = args.Get(0).(*models.AlbumWithArtist)
	}

	return album, args.Error(1)
}

func (m *MockAlbumRepository) FindAlbumsByArtistId(ctx context.Context, artistId int) (albums []models.Album, err error) {
	args := m.Called(ctx, artistId)

	if args.Get(0) != nil {
		albums = args.Get(0).([]models.Album)
	}

	return albums, args.Error(1)
}

func (m *MockAlbumRepository) FindAll(ctx context.Context, pageSize int, offset int) (albums []models.Album, err error) {
	args := m.Called(ctx, pageSize, offset)

	if args.Get(0) != nil {
		albums = args.Get(0).([]models.Album)
	}

	return albums, args.Error(1)
}

func (m *MockAlbumRepository) FindCount(ctx context.Context) (total int, err error) {
	args := m.Called(ctx)

	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *MockAlbumRepository) FindExistsAlbumById(ctx context.Context, id int) (exists bool, err error) {
	args := m.Called(ctx, id)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockAlbumRepository) FindExistsAlbumBySlug(ctx context.Context, slug string) (exists bool, err error) {
	args := m.Called(ctx, slug)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockAlbumRepository) Store(ctx context.Context, input models.CreateAlbumInput) (err error) {
	args := m.Called(ctx, input)

	return args.Error(0)
}

func (m *MockAlbumRepository) Update(ctx context.Context, input models.CreateAlbumInput, id int) (err error) {
	args := m.Called(ctx, input, id)

	return args.Error(0)
}
