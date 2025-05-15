package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/mocks"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type AlbumServiceTestSuite struct {
	suite.Suite
	Svc        contracts.AlbumService
	AlbumRepo  *mocks.MockAlbumRepository
	ArtistRepo *mocks.MockArtistRepository
}

func (s *AlbumServiceTestSuite) SetupTest() {
	s.AlbumRepo = new(mocks.MockAlbumRepository)
	s.ArtistRepo = new(mocks.MockArtistRepository)
	s.Svc = NewAlbumService(s.AlbumRepo, s.ArtistRepo, nil)
}

func (s *AlbumServiceTestSuite) ResetMocks() {
	s.AlbumRepo.Calls = nil
	s.AlbumRepo.ExpectedCalls = nil
	s.ArtistRepo.Calls = nil
	s.ArtistRepo.ExpectedCalls = nil
}

func (s *AlbumServiceTestSuite) TestGetAll() {
	image := dto.Image{Src: "image.png", BlurHash: "abcs"}

	testCases := []struct {
		name            string
		prepareMock     func()
		expectedResults []dto.AlbumWithArtist
		expectedTotal   int
		expectErr       error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.AlbumRepo.On("FindCount", mock.Anything).Return(1, nil)
				s.AlbumRepo.On("FindAll", mock.Anything, 10, 0).Return([]models.Album{
					{
						Id:       1,
						ArtistId: 1,
						Name:     "Bintang di surga",
						Slug:     "bintang-di-surga",
						Image:    utils.ParseImageToByte(&image),
					},
				}, nil)

				artistIds := []any{1}
				inClause, args := utils.BuildInClause(1, artistIds)

				s.ArtistRepo.On("FindByArtistIds", mock.Anything, inClause, args).Return([]models.Artist{
					{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
			},
			expectedResults: []dto.AlbumWithArtist{
				{
					Album: dto.Album{
						Id:    1,
						Name:  "Bintang di surga",
						Slug:  "bintang-di-surga",
						Image: image,
					},
					Artist: dto.Artist{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: image,
					},
				},
			},
			expectedTotal: 1,
		},
		{
			name: "FindCount_Error",
			prepareMock: func() {
				s.AlbumRepo.On("FindCount", mock.Anything).Return(0, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "FindAll_Error",
			prepareMock: func() {
				s.AlbumRepo.On("FindCount", mock.Anything).Return(1, nil)
				s.AlbumRepo.On("FindAll", mock.Anything, 10, 0).Return(nil, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "FindByArtistIds_Error",
			prepareMock: func() {
				s.AlbumRepo.On("FindCount", mock.Anything).Return(1, nil)
				s.AlbumRepo.On("FindAll", mock.Anything, 10, 0).Return([]models.Album{
					{
						Id:       1,
						ArtistId: 1,
						Name:     "Bintang di surga",
						Slug:     "bintang-di-surga",
						Image:    utils.ParseImageToByte(&image),
					},
				}, nil)

				artistIds := []any{1}
				inClause, args := utils.BuildInClause(1, artistIds)

				s.ArtistRepo.On("FindByArtistIds", mock.Anything, inClause, args).Return(nil, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			// Actual
			results, total, err := s.Svc.GetAll(s.T().Context(), 10, 0)

			// Assertion
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectedResults, results)
				s.Equal(tc.expectedTotal, total)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.AlbumRepo.AssertExpectations(s.T())
			s.ArtistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AlbumServiceTestSuite) TestGetAlbumById() {
	image := dto.Image{Src: "image.png", BlurHash: "abcs"}

	testCases := []struct {
		name           string
		prepareMock    func()
		expectedResult dto.AlbumWithArtist
		expectErr      error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumById", mock.Anything, 1).Return(&models.AlbumWithArtist{
					Album: models.Album{
						Id:       1,
						ArtistId: 1,
						Name:     "Bintang di surga",
						Slug:     "bintang-di-surga",
						Image:    utils.ParseImageToByte(&image),
					},
					Artist: models.Artist{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
			},
			expectedResult: dto.AlbumWithArtist{
				Album: dto.Album{
					Id:    1,
					Name:  "Bintang di surga",
					Slug:  "bintang-di-surga",
					Image: image,
				},
				Artist: dto.Artist{
					Id:    1,
					Name:  "Noah",
					Slug:  "noah",
					Image: image,
				},
			},
		},
		{
			name: "FindAlbumById_Error",
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumById", mock.Anything, 1).Return(nil, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			// Actual
			result, err := s.Svc.GetAlbumById(s.T().Context(), 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectedResult, result)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.AlbumRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AlbumServiceTestSuite) TestCreateAlbum() {
	image := dto.Image{Src: "image.png", BlurHash: "abc"}
	slug := utils.MakeSlug("Test Album")
	var validationErr = errors.New("validation failed")

	testCases := []struct {
		name               string
		createAlbumRequest dto.CreateAlbumRequest
		prepareMock        func()
		expectErr          error
		expectValErrorMap  map[string]string
	}{
		{
			name: "success",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("FindExistsAlbumBySlug", mock.Anything, slug).Return(false, nil)
				s.AlbumRepo.On("Store", mock.Anything, models.CreateAlbumInput{
					Name:     "Test Album",
					ArtistId: 1,
					Slug:     slug,
					Image:    utils.ParseImageToByte(&image),
				}).Return(nil)
			},
		},
		{
			name: "ValidationErrors_RequiredAll",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "",
				ArtistId: 0,
				Image:    nil,
			},
			expectErr: validationErr,
			expectValErrorMap: map[string]string{
				"name":      "Field is required",
				"artist_id": "Field is required",
				"image":     "Field is required",
			},
		},
		{
			name: "ValidationErrors_RequiredName",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "",
				ArtistId: 1,
				Image:    &image,
			},
			expectErr:         validationErr,
			expectValErrorMap: map[string]string{"name": "Field is required"},
		},
		{
			name: "ValidationErrors_RequiredArtistId",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 0,
				Image:    &image,
			},
			expectErr:         validationErr,
			expectValErrorMap: map[string]string{"artist_id": "Field is required"},
		},
		{
			name: "ValidationErrors_RequiredImage",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    nil,
			},
			expectErr:         validationErr,
			expectValErrorMap: map[string]string{"image": "Field is required"},
		},
		{
			name: "FindExistsArtistById_NotFound",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Artist", "id", 1),
		},
		{
			name: "FindExistsAlbumBySlug_Conflict",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("FindExistsAlbumBySlug", mock.Anything, slug).Return(true, nil)
			},
			expectErr: errs.NewConflictError("Album", "name", "Test Album"),
		},
		{
			name: "Store_Error",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("FindExistsAlbumBySlug", mock.Anything, slug).Return(false, nil)
				s.AlbumRepo.On("Store", mock.Anything, models.CreateAlbumInput{
					Name:     "Test Album",
					ArtistId: 1,
					Slug:     slug,
					Image:    utils.ParseImageToByte(&image),
				}).Return(errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			// Actual
			err := s.Svc.CreateAlbum(s.T().Context(), tc.createAlbumRequest)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else if errors.Is(tc.expectErr, validationErr) {
				var valErr *errs.BadRequestError
				s.ErrorAs(err, &valErr)
				s.Equal(tc.expectValErrorMap, valErr.Errors)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.ArtistRepo.AssertExpectations(s.T())
			s.AlbumRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AlbumServiceTestSuite) TestUpdateAlbum() {
	image := dto.Image{Src: "image.png", BlurHash: "abc"}
	slug := utils.MakeSlug("Test Album")
	var validationErr = errors.New("validation failed")

	testCases := []struct {
		name               string
		createAlbumRequest dto.CreateAlbumRequest
		prepareMock        func()
		expectErr          error
		expectValErrorMap  map[string]string
	}{
		{
			name: "success_withoutDataChanges", // Scenario success1 for Update album without data changes.
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumById", mock.Anything, 1).Return(&models.AlbumWithArtist{
					Album: models.Album{
						Id:       1,
						ArtistId: 1,
						Name:     "Test Album",
						Slug:     "test-album",
						Image:    utils.ParseImageToByte(&image),
					},
					Artist: models.Artist{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("Update", mock.Anything, models.CreateAlbumInput{
					Name:     "Test Album",
					ArtistId: 1,
					Slug:     slug,
					Image:    utils.ParseImageToByte(&image),
				}, 1).Return(nil)
			},
		},
		{
			name: "success_withDataChange", // Scenario success2 for Update album with data changes.
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Nama Album Baru",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumById", mock.Anything, 1).Return(&models.AlbumWithArtist{
					Album: models.Album{
						Id:       1,
						ArtistId: 1,
						Name:     "Test Album",
						Slug:     "test-album",
						Image:    utils.ParseImageToByte(&image),
					},
					Artist: models.Artist{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("FindExistsAlbumBySlug", mock.Anything, utils.MakeSlug("Nama Album Baru")).Return(false, nil)
				s.AlbumRepo.On("Update", mock.Anything, models.CreateAlbumInput{
					Name:     "Nama Album Baru",
					ArtistId: 1,
					Slug:     utils.MakeSlug("Nama Album Baru"),
					Image:    utils.ParseImageToByte(&image),
				}, 1).Return(nil)
			},
		},
		{
			name: "ValidationErrors_RequiredAll",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "",
				ArtistId: 0,
				Image:    nil,
			},
			expectErr: validationErr,
			expectValErrorMap: map[string]string{
				"name":      "Field is required",
				"artist_id": "Field is required",
				"image":     "Field is required",
			},
		},
		{
			name: "ValidationErrors_RequiredName",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "",
				ArtistId: 1,
				Image:    &image,
			},
			expectErr:         validationErr,
			expectValErrorMap: map[string]string{"name": "Field is required"},
		},
		{
			name: "ValidationErrors_RequiredArtistId",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 0,
				Image:    &image,
			},
			expectErr:         validationErr,
			expectValErrorMap: map[string]string{"artist_id": "Field is required"},
		},
		{
			name: "ValidationErrors_RequiredImage",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    nil,
			},
			expectErr:         validationErr,
			expectValErrorMap: map[string]string{"image": "Field is required"},
		},
		{
			name: "FindExistsArtistById_NotFound",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Test Album",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumById", mock.Anything, 1).Return(&models.AlbumWithArtist{
					Album: models.Album{
						Id:       1,
						ArtistId: 1,
						Name:     "Test Album",
						Slug:     "test-album",
						Image:    utils.ParseImageToByte(&image),
					},
					Artist: models.Artist{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Artist", "id", 1),
		},
		{
			name: "FindExistsAlbumBySlug_ConflictWhileAlbumNameHasChange",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Nama Album Baru",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumById", mock.Anything, 1).Return(&models.AlbumWithArtist{
					Album: models.Album{
						Id:       1,
						ArtistId: 1,
						Name:     "Test Album",
						Slug:     "test-album",
						Image:    utils.ParseImageToByte(&image),
					},
					Artist: models.Artist{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("FindExistsAlbumBySlug", mock.Anything, utils.MakeSlug("Nama Album Baru")).Return(true, nil)
			},
			expectErr: errs.NewConflictError("Album", "name", "Nama Album Baru"),
		},
		{
			name: "Update_Error",
			createAlbumRequest: dto.CreateAlbumRequest{
				Name:     "Nama Album Baru",
				ArtistId: 1,
				Image:    &image,
			},
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumById", mock.Anything, 1).Return(&models.AlbumWithArtist{
					Album: models.Album{
						Id:       1,
						ArtistId: 1,
						Name:     "Test Album",
						Slug:     "test-album",
						Image:    utils.ParseImageToByte(&image),
					},
					Artist: models.Artist{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("FindExistsAlbumBySlug", mock.Anything, utils.MakeSlug("Nama Album Baru")).Return(false, nil)
				s.AlbumRepo.On("Update", mock.Anything, models.CreateAlbumInput{
					Name:     "Nama Album Baru",
					ArtistId: 1,
					Slug:     utils.MakeSlug("Nama Album Baru"),
					Image:    utils.ParseImageToByte(&image),
				}, 1).Return(errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			// Actual
			err := s.Svc.UpdateAlbum(s.T().Context(), tc.createAlbumRequest, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else if errors.Is(tc.expectErr, validationErr) {
				var valErr *errs.BadRequestError
				s.ErrorAs(err, &valErr)
				s.Equal(tc.expectValErrorMap, valErr.Errors)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.ArtistRepo.AssertExpectations(s.T())
			s.AlbumRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AlbumServiceTestSuite) TestDeleteAlbum() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.AlbumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("Delete", mock.Anything, 1).Return(nil)
			},
		},
		{
			name: "FindExistsAlbumById_NotFound",
			prepareMock: func() {
				s.AlbumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Album", "id", 1),
		},
		{
			name: "FindExistsAlbumById_Error",
			prepareMock: func() {
				s.AlbumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(false, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "Delete_Error",
			prepareMock: func() {
				s.AlbumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(true, nil)
				s.AlbumRepo.On("Delete", mock.Anything, 1).Return(errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			// Actual
			err := s.Svc.DeleteAlbum(s.T().Context(), 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.AlbumRepo.AssertExpectations(s.T())
		})
	}
}

func (s *AlbumServiceTestSuite) TestGetAlbumsByArtistId() {
	image := dto.Image{Src: "image.png", BlurHash: "abcd"}
	testCases := []struct {
		name          string
		prepareMock   func()
		expectResults []dto.Album
		expectErr     error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumsByArtistId", mock.Anything, 1).Return([]models.Album{
					{
						Id:       1,
						ArtistId: 1,
						Name:     "Test Album",
						Slug:     "test-album",
						Image:    utils.ParseImageToByte(&image),
					},
				}, nil)
			},
			expectResults: []dto.Album{
				{
					Id:    1,
					Name:  "Test Album",
					Slug:  "test-album",
					Image: image,
				},
			},
		},
		{
			name: "FindAlbumsByArtistId_Error",
			prepareMock: func() {
				s.AlbumRepo.On("FindAlbumsByArtistId", mock.Anything, 1).Return(nil, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			if tc.prepareMock != nil {
				tc.prepareMock()
			}

			// Actual
			results, err := s.Svc.GetAlbumsByArtistId(s.T().Context(), 1)

			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectResults, results)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.AlbumRepo.AssertExpectations(s.T())
		})
	}
}

func TestAlbumServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AlbumServiceTestSuite))
}
