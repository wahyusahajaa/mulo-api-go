package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/mocks"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
)

type GenreServiceTestSuite struct {
	suite.Suite
	Svc            contracts.GenreService
	MockGenreRepo  *mocks.MockGenreRepository
	MockArtistRepo *mocks.MockArtistRepository
	MockSongRepo   *mocks.MockSongRepository
}

func (s *GenreServiceTestSuite) SetupTest() {
	s.MockGenreRepo = new(mocks.MockGenreRepository)
	s.MockArtistRepo = new(mocks.MockArtistRepository)
	s.MockSongRepo = new(mocks.MockSongRepository)
	s.Svc = NewGenreService(s.MockGenreRepo, s.MockArtistRepo, s.MockSongRepo, nil)
}

func (s *GenreServiceTestSuite) ResetMocks() {
	s.MockGenreRepo.ExpectedCalls = nil
	s.MockGenreRepo.Calls = nil
	s.MockArtistRepo.ExpectedCalls = nil
	s.MockArtistRepo.Calls = nil
	s.MockSongRepo.ExpectedCalls = nil
	s.MockSongRepo.Calls = nil
}

func (s *GenreServiceTestSuite) TestGetAll() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abcd"}
	image2 := dto.Image{Src: "image2.png", BlurHash: "cdsa"}
	image1Bytes, _ := json.Marshal(image1)
	image2Bytes, _ := json.Marshal(image2)

	type expected struct {
		results []dto.Genre
		total   int
		err     error
	}
	type scenario struct {
		name        string
		prepareMock func()
		expected    expected
	}

	testCases := []scenario{
		{
			name: "success",
			prepareMock: func() {
				s.MockGenreRepo.On("FindCount", mock.Anything).Return(2, nil)
				s.MockGenreRepo.On("FindAll", mock.Anything, 10, 0).Return([]models.Genre{
					{Id: 1, Name: "Genre 1", Image: image1Bytes},
					{Id: 2, Name: "Genre 2", Image: image2Bytes},
				}, nil)
			},
			expected: expected{
				results: []dto.Genre{{Id: 1, Name: "Genre 1", Image: image1}, {Id: 2, Name: "Genre 2", Image: image2}},
				total:   2,
			},
		},
		{
			name: "findCountError",
			prepareMock: func() {
				s.MockGenreRepo.On("FindCount", mock.Anything).Return(0, errors.New("database failure"))
			},
			expected: expected{
				err: errors.New("database failure"),
			},
		},
		{
			name: "findAllError",
			prepareMock: func() {
				s.MockGenreRepo.On("FindCount", mock.Anything).Return(2, nil)
				s.MockGenreRepo.On("FindAll", mock.Anything, 10, 0).Return(nil, errors.New("database failure"))
			},
			expected: expected{
				err: errors.New("database failure"),
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			results, total, err := s.Svc.GetAll(s.T().Context(), 10, 0)

			if tc.expected.err == nil {
				s.NoError(err)
				s.Equal(tc.expected.results, results)
				s.Equal(tc.expected.total, total)
			} else {
				s.Error(err)
				s.EqualError(tc.expected.err, err.Error())
			}

			s.MockGenreRepo.AssertExpectations(s.T())
		})
	}

}

func (s *GenreServiceTestSuite) TestGetGenreById() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abcd"}
	image1Bytes, _ := json.Marshal(image1)

	type scenario struct {
		name           string
		prepareMock    func()
		expectedResult dto.Genre
		expectedErr    error
	}

	testCases := []scenario{
		{
			name: "success",
			prepareMock: func() {
				s.MockGenreRepo.On("FindGenreById", mock.Anything, 1).Return(&models.Genre{
					Id:    1,
					Name:  "Genre 1",
					Image: image1Bytes,
				}, nil)
			},
			expectedResult: dto.Genre{Id: 1, Name: "Genre 1", Image: image1},
			expectedErr:    nil,
		},
		{
			name: "FindGenreById_NotFound",
			prepareMock: func() {
				s.MockGenreRepo.On("FindGenreById", mock.Anything, 1).Return(nil, nil)
			},
			expectedErr: errs.NewNotFoundError("Genre", "id", 1),
		},
		{
			name: "FindGenreById_Error",
			prepareMock: func() {
				s.MockGenreRepo.On("FindGenreById", mock.Anything, 1).Return(nil, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			result, err := s.Svc.GetGenreById(s.T().Context(), 1)

			if tc.expectedErr == nil {
				s.NoError(err)
				s.Equal(tc.expectedResult, result)
			} else {
				s.Error(err)
				s.EqualError(tc.expectedErr, err.Error())
			}

			s.MockGenreRepo.AssertExpectations(s.T())
		})
	}
}

func (s *GenreServiceTestSuite) TestCreateGenre() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abcd"}
	image1Bytes, _ := json.Marshal(image1)
	var errValidation = errors.New("validation failed")

	type scenario struct {
		name              string
		prepareMock       func()
		reqDto            dto.CreateGenreRequest
		expectedErr       error
		expectedValErrMap map[string]string
	}

	testCases := []scenario{
		{
			name: "success",
			prepareMock: func() {
				s.MockGenreRepo.On("Store", mock.Anything, models.CreateGenreInput{
					Name:  "Genre 1",
					Image: image1Bytes,
				}).Return(nil)
			},
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: &image1,
			},
			expectedErr: nil,
		},
		{
			name: "storeError",
			prepareMock: func() {
				s.MockGenreRepo.On("Store", mock.Anything, models.CreateGenreInput{
					Name:  "Genre 1",
					Image: image1Bytes,
				}).Return(errors.New("database failure"))
			},
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: &image1,
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "validationErrors_RequiredName",
			reqDto: dto.CreateGenreRequest{
				Name:  "",
				Image: &image1,
			},
			expectedErr:       errValidation,
			expectedValErrMap: map[string]string{"name": "Field is required"},
		},
		{
			name: "validationErrors_RequiredImage",
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: nil,
			},
			expectedErr:       errValidation,
			expectedValErrMap: map[string]string{"image": "Field is required"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()

			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			err := s.Svc.CreateGenre(context.Background(), tc.reqDto)

			if tc.expectedErr == nil {
				s.NoError(err)
			} else if errors.Is(tc.expectedErr, errValidation) {
				var valErr *errs.BadRequestError
				s.ErrorAs(err, &valErr)
				s.Equal(tc.expectedValErrMap, valErr.Errors)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.MockGenreRepo.AssertExpectations(s.T())
		})
	}

}

func (s *GenreServiceTestSuite) TestDeleteGenre() {
	type expected struct {
		err error
	}

	type scenario struct {
		name        string
		prepareMock func()
		expected    expected
	}

	testCases := []scenario{
		{
			name: "success",
			prepareMock: func() {
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("Delete", mock.Anything, 1).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "genreNotFound",
			prepareMock: func() {
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(false, errs.NewNotFoundError("Genre", "id", 1))
			},
			expected: expected{
				err: errs.NewNotFoundError("Genre", "id", 1),
			},
		},
		{
			name: "deleteError",
			prepareMock: func() {
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("Delete", mock.Anything, 1).Return(errors.New("database failure"))
			},
			expected: expected{
				err: errors.New("database failure"),
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			// Actual
			err := s.Svc.DeleteGenre(context.Background(), 1)

			// Assert
			if tc.expected.err == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expected.err.Error())
			}

			s.MockGenreRepo.AssertExpectations(s.T())
		})
	}
}

func (s *GenreServiceTestSuite) TestUpdateGenre() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abc"}
	image1Bytes, _ := json.Marshal(image1)
	var errValidation = errors.New("validation failed")

	type scenario struct {
		name                string
		reqDto              dto.CreateGenreRequest
		prepareMock         func()
		expectedErr         error
		expectedValErrorMap map[string]string
	}

	testCases := []scenario{
		{
			name: "success",
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: &image1,
			},
			prepareMock: func() {
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("Update", mock.Anything, models.CreateGenreInput{
					Name:  "Genre 1",
					Image: image1Bytes,
				}, 1).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "validationErrors_MissingName",
			reqDto: dto.CreateGenreRequest{
				Name:  "",
				Image: &image1,
			},
			expectedErr:         errValidation,
			expectedValErrorMap: map[string]string{"name": "Field is required"},
		},
		{
			name: "validationErrors_MissingImage",
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: nil,
			},
			expectedErr:         errValidation,
			expectedValErrorMap: map[string]string{"image": "Field is required"},
		},
		{
			name: "FindExistsGenreById_Error",
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: &image1,
			},
			prepareMock: func() {
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(false, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "FindExistsGenreById_NotFound",
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: &image1,
			},
			prepareMock: func() {
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(false, nil)
			},
			expectedErr: errs.NewNotFoundError("Genre", "Id", 1),
		},
		{
			name: "updateError",
			reqDto: dto.CreateGenreRequest{
				Name:  "Genre 1",
				Image: &image1,
			},
			prepareMock: func() {
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("Update", mock.Anything, models.CreateGenreInput{
					Name:  "Genre 1",
					Image: image1Bytes,
				}, 1).Return(errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			err := s.Svc.UpdateGenre(s.T().Context(), tc.reqDto, 1)

			if tc.expectedErr == nil {
				s.NoError(err)
			} else if errors.Is(tc.expectedErr, errValidation) {
				var valErr *errs.BadRequestError
				s.ErrorAs(err, &valErr)
				s.Equal(tc.expectedValErrorMap, valErr.Errors)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.MockGenreRepo.AssertExpectations(s.T())
		})
	}

}

func (s *GenreServiceTestSuite) TestCreateArtistGenre() {
	type scenario struct {
		name        string
		prepareMock func()
		expectedErr error
	}

	testCases := []scenario{
		{
			name: "success",
			prepareMock: func() {
				s.MockArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("FindExistsArtistGenreByGenreId", mock.Anything, 1, 1).Return(false, nil)
				s.MockGenreRepo.On("StoreArtistGenre", mock.Anything, 1, 1).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "FindExistsArtistById_NotFound",
			prepareMock: func() {
				s.MockArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(false, nil)
			},
			expectedErr: errs.NewNotFoundError("Artist", "id", 1),
		},
		{
			name: "FindExistsGenreById_NotFound",
			prepareMock: func() {
				s.MockArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(false, nil)
			},
			expectedErr: errs.NewNotFoundError("Genre", "id", 1),
		},
		{
			name: "FindExistsArtistGenreByGenreId_Conflict",
			prepareMock: func() {
				s.MockArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("FindExistsArtistGenreByGenreId", mock.Anything, 1, 1).Return(true, nil)
			},
			expectedErr: errs.NewConflictError("Genre", "genre_id", 1),
		},
		{
			name: "StoreArtistGenre_Error",
			prepareMock: func() {
				s.MockArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("FindExistsGenreById", mock.Anything, 1).Return(true, nil)
				s.MockGenreRepo.On("FindExistsArtistGenreByGenreId", mock.Anything, 1, 1).Return(false, nil)
				s.MockGenreRepo.On("StoreArtistGenre", mock.Anything, 1, 1).Return(errors.New("db failure"))
			},
			expectedErr: errors.New("db failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			err := s.Svc.CreateArtistGenre(s.T().Context(), 1, 1)

			if tc.expectedErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.MockArtistRepo.AssertExpectations(s.T())
			s.MockGenreRepo.AssertExpectations(s.T())
		})
	}
}

func (s *GenreServiceTestSuite) TestGetArtistGenres() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abcd"}
	image2 := dto.Image{Src: "image2.png", BlurHash: "bcda"}
	image1Bytes, _ := json.Marshal(image1)
	image2Bytes, _ := json.Marshal(image2)

	type scenario struct {
		name            string
		prepareMock     func()
		expectedResults []dto.Genre
		expectedErr     error
	}

	testCases := []scenario{
		{
			name: "success",
			prepareMock: func() {
				s.MockGenreRepo.On("FindArtistGenres", mock.Anything, 1, 10, 0).Return([]models.Genre{
					{Id: 1, Name: "Genre 1", Image: image1Bytes},
					{Id: 2, Name: "Genre 2", Image: image2Bytes},
				}, nil)
			},
			expectedResults: []dto.Genre{
				{Id: 1, Name: "Genre 1", Image: image1},
				{Id: 2, Name: "Genre 2", Image: image2},
			},
			expectedErr: nil,
		},
		{
			name: "error",
			prepareMock: func() {
				s.MockGenreRepo.On("FindArtistGenres", mock.Anything, 1, 10, 0).Return(nil, errors.New("db failure"))
			},
			expectedErr: errors.New("db failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			results, err := s.Svc.GetArtistGenres(s.T().Context(), 1, 10, 0)

			if tc.expectedErr == nil {
				s.NoError(err)
				s.Equal(tc.expectedResults, results)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.MockGenreRepo.AssertExpectations(s.T())
		})
	}
}

func TestGenreServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GenreServiceTestSuite))
}
