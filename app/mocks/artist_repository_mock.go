package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type MockArtistRepository struct {
	mock.Mock
}

func (m *MockArtistRepository) Delete(ctx context.Context, id int) (err error) {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *MockArtistRepository) FindAll(ctx context.Context, pageSize int, offset int) (artists []models.Artist, err error) {
	args := m.Called(ctx, pageSize, offset)

	if args.Get(0) != nil {
		artists = args.Get(0).([]models.Artist)
	}

	return artists, args.Error(1)
}

func (m *MockArtistRepository) FindArtistById(ctx context.Context, artistId int) (artist *models.Artist, err error) {
	args := m.Called(ctx, artistId)

	if args.Get(0) != nil {
		artist = args.Get(0).(*models.Artist)
	}

	return artist, args.Error(1)
}

func (m *MockArtistRepository) FindByArtistIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error) {
	args := m.Called(ctx, inClause, artistIds)

	if args.Get(0) != nil {
		artists = args.Get(0).([]models.Artist)
	}

	return artists, args.Error(1)
}

func (m *MockArtistRepository) FindCount(ctx context.Context) (total int, err error) {
	args := m.Called(ctx)

	if args.Get(0) != nil {
		total = args.Get(0).(int)
	}

	return total, args.Error(1)
}

func (m *MockArtistRepository) FindExistsArtistBySlug(ctx context.Context, slug string) (exists bool, err error) {
	args := m.Called(ctx, slug)

	if args.Get(0) != nil {
		exists = args.Get(0).(bool)
	}

	return exists, args.Error(1)
}

func (m *MockArtistRepository) Store(ctx context.Context, input models.CreateArtistInput) (err error) {
	args := m.Called(ctx, input)

	return args.Error(0)
}

func (m *MockArtistRepository) Update(ctx context.Context, input models.CreateArtistInput, id int) (err error) {
	args := m.Called(ctx, input, id)

	return args.Error(0)
}

func (m *MockArtistRepository) FindExistsArtistById(ctx context.Context, id int) (exists bool, err error) {
	args := m.Called(ctx, id)

	return args.Get(0).(bool), args.Error(1)
}
