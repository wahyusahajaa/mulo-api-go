package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type MockArtistRepository struct {
	mock.Mock
}

// Delete implements contracts.ArtistRepository.
func (m *MockArtistRepository) Delete(ctx context.Context, id int) (err error) {
	panic("unimplemented")
}

// FindAll implements contracts.ArtistRepository.
func (m *MockArtistRepository) FindAll(ctx context.Context, pageSize int, offset int) (artists []models.Artist, err error) {
	panic("unimplemented")
}

// FindArtistById implements contracts.ArtistRepository.
func (m *MockArtistRepository) FindArtistById(ctx context.Context, artistId int) (artist *models.Artist, err error) {
	panic("unimplemented")
}

// FindByArtistIds implements contracts.ArtistRepository.
func (m *MockArtistRepository) FindByArtistIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error) {
	panic("unimplemented")
}

// FindCount implements contracts.ArtistRepository.
func (m *MockArtistRepository) FindCount(ctx context.Context) (total int, err error) {
	panic("unimplemented")
}

// FindExistsArtistBySlug implements contracts.ArtistRepository.
func (m *MockArtistRepository) FindExistsArtistBySlug(ctx context.Context, slug string) (exists bool, err error) {
	panic("unimplemented")
}

// Store implements contracts.ArtistRepository.
func (m *MockArtistRepository) Store(ctx context.Context, input models.CreateArtistInput) (err error) {
	panic("unimplemented")
}

// Update implements contracts.ArtistRepository.
func (m *MockArtistRepository) Update(ctx context.Context, input models.CreateArtistInput, id int) (err error) {
	panic("unimplemented")
}

func (m *MockArtistRepository) FindExistsArtistById(ctx context.Context, id int) (exists bool, err error) {
	args := m.Called(ctx, id)

	return args.Get(0).(bool), args.Error(1)
}
