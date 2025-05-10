package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
)

type MockGenreService struct {
	mock.Mock
}

func (m *MockGenreService) GetAll(ctx context.Context, pageSize, offset int) (genres []dto.Genre, err error) {
	args := m.Called(ctx, pageSize, offset)

	if args.Get(0) != nil {
		genres = args.Get(0).([]dto.Genre)
	}
	err = args.Error(1)
	return
}
