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

type ArtistServiceTestSuite struct {
	suite.Suite
	Svc        contracts.ArtistService
	ArtistRepo *mocks.MockArtistRepository
}

func (s *ArtistServiceTestSuite) SetupTest() {
	s.ArtistRepo = new(mocks.MockArtistRepository)
	s.Svc = NewArtistService(s.ArtistRepo, nil)
}

func (s *ArtistServiceTestSuite) ResetMocks() {
	s.ArtistRepo.ExpectedCalls = nil
	s.ArtistRepo.Calls = nil
}

func (s *ArtistServiceTestSuite) TestGetAll() {
	image := dto.Image{Src: "image1.png", BlurHash: "abc"}
	testCases := []struct {
		name          string
		prepareMock   func()
		expectResults []dto.Artist
		expectTotal   int
		expectErr     error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.ArtistRepo.On("FindCount", mock.Anything).Return(1, nil)
				s.ArtistRepo.On("FindAll", mock.Anything, 10, 0).Return([]models.Artist{
					{
						Id:    1,
						Name:  "Noah",
						Slug:  "noah",
						Image: utils.ParseImageToByte(&image),
					},
				}, nil)
			},
			expectResults: []dto.Artist{
				{
					Id:    1,
					Name:  "Noah",
					Slug:  "noah",
					Image: image,
				},
			},
			expectTotal: 1,
		},
		{
			name: "FindCount_Error",
			prepareMock: func() {
				s.ArtistRepo.On("FindCount", mock.Anything).Return(0, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "FindCount_Error",
			prepareMock: func() {
				s.ArtistRepo.On("FindCount", mock.Anything).Return(1, nil)
				s.ArtistRepo.On("FindAll", mock.Anything, 10, 0).Return(nil, errors.New("database failure"))
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

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectResults, results)
				s.Equal(tc.expectTotal, total)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.ArtistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *ArtistServiceTestSuite) TestGetArtistById() {
	image := dto.Image{Src: "image.png", BlurHash: "abc"}
	testCases := []struct {
		name         string
		prepareMock  func()
		expectResult dto.Artist
		expectErr    error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(&models.Artist{
					Id:    1,
					Name:  "Noah",
					Slug:  "noah",
					Image: utils.ParseImageToByte(&image),
				}, nil)
			},
			expectResult: dto.Artist{
				Id:    1,
				Name:  "Noah",
				Slug:  "noah",
				Image: image,
			},
		},
		{
			name: "GetArtistById_NotFound",
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(nil, nil)
			},
			expectErr: errs.NewNotFoundError("Artist", "id", 1),
		},
		{
			name: "FindArtistById_Error",
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(&models.Artist{
					Id:    1,
					Name:  "Noah",
					Slug:  "noah",
					Image: utils.ParseImageToByte(&image),
				}, errors.New("database failure"))
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
			result, err := s.Svc.GetArtistById(s.T().Context(), 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectResult, result)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.ArtistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *ArtistServiceTestSuite) TestCreateArtist() {
	image := dto.Image{Src: "image.png", BlurHash: "abc"}
	artistName := "Ungu"
	var validationErr = errors.New("validation failed")

	testCases := []struct {
		name            string
		req             dto.CreateArtistRequest
		prepareMock     func()
		expectErr       error
		expectValErrMap map[string]string
	}{
		{
			name: "success",
			req: dto.CreateArtistRequest{
				Name:  artistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistBySlug", mock.Anything, utils.MakeSlug(artistName)).Return(false, nil)
				s.ArtistRepo.On("Store", mock.Anything, models.CreateArtistInput{
					Name:  artistName,
					Slug:  utils.MakeSlug(artistName),
					Image: utils.ParseImageToByte(&image),
				}).Return(nil)
			},
		},
		{
			name: "validationErrors_RequiredName",
			req: dto.CreateArtistRequest{
				Name:  "",
				Image: &image,
			},
			expectErr:       validationErr,
			expectValErrMap: map[string]string{"name": "Field is required"},
		},
		{
			name: "validationErrors_RequiredImage",
			req: dto.CreateArtistRequest{
				Name:  "Noah",
				Image: nil,
			},
			expectErr:       validationErr,
			expectValErrMap: map[string]string{"image": "Field is required"},
		},
		{
			name: "FindExistsArtistBySlug_Conflict",
			req: dto.CreateArtistRequest{
				Name:  artistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistBySlug", mock.Anything, utils.MakeSlug(artistName)).Return(true, nil)
			},
			expectErr: errs.NewConflictError("Artist", "name", artistName),
		},
		{
			name: "FindExistsArtistBySlug_Error",
			req: dto.CreateArtistRequest{
				Name:  artistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistBySlug", mock.Anything, utils.MakeSlug(artistName)).Return(false, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "Store_Error",
			req: dto.CreateArtistRequest{
				Name:  artistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistBySlug", mock.Anything, utils.MakeSlug(artistName)).Return(false, nil)
				s.ArtistRepo.On("Store", mock.Anything, models.CreateArtistInput{
					Name:  artistName,
					Slug:  utils.MakeSlug(artistName),
					Image: utils.ParseImageToByte(&image),
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
			err := s.Svc.CreateArtist(s.T().Context(), tc.req)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else if errors.Is(tc.expectErr, validationErr) {
				var valErr *errs.BadRequestError
				s.ErrorAs(err, &valErr)
				s.Equal(tc.expectValErrMap, valErr.Errors)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.ArtistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *ArtistServiceTestSuite) TestUpdateArtist() {
	image := dto.Image{Src: "image.png", BlurHash: "abc"}
	oldArtistName := "Noah"
	newArtistName := "Ungu"

	testCases := []struct {
		name        string
		req         dto.CreateArtistRequest
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success_withoutArtistChanges",
			req: dto.CreateArtistRequest{
				Name:  oldArtistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(&models.Artist{
					Id:    1,
					Name:  oldArtistName,
					Slug:  utils.MakeSlug(oldArtistName),
					Image: utils.ParseImageToByte(&image),
				}, nil)

				s.ArtistRepo.On("Update", mock.Anything, models.CreateArtistInput{
					Name:  oldArtistName,
					Slug:  utils.MakeSlug(oldArtistName),
					Image: utils.ParseImageToByte(&image),
				}, 1).Return(nil)
			},
		},
		{
			name: "success_withArtistChanges",
			req: dto.CreateArtistRequest{
				Name:  newArtistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(&models.Artist{
					Id:    1,
					Name:  oldArtistName,
					Slug:  utils.MakeSlug(oldArtistName),
					Image: utils.ParseImageToByte(&image),
				}, nil)

				s.ArtistRepo.On("FindExistsArtistBySlug", mock.Anything, utils.MakeSlug(newArtistName)).Return(false, nil)

				s.ArtistRepo.On("Update", mock.Anything, models.CreateArtistInput{
					Name:  newArtistName,
					Slug:  utils.MakeSlug(newArtistName),
					Image: utils.ParseImageToByte(&image),
				}, 1).Return(nil)
			},
		},
		{
			name: "FindArtistById_NotFound",
			req: dto.CreateArtistRequest{
				Name:  newArtistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(nil, nil)
			},
			expectErr: errs.NewNotFoundError("Artist", "id", 1),
		},
		{
			name: "FindExistsArtistBySlug_Conflict_withArtistChanges",
			req: dto.CreateArtistRequest{
				Name:  newArtistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(&models.Artist{
					Id:    1,
					Name:  oldArtistName,
					Slug:  utils.MakeSlug(oldArtistName),
					Image: utils.ParseImageToByte(&image),
				}, nil)

				s.ArtistRepo.On("FindExistsArtistBySlug", mock.Anything, utils.MakeSlug(newArtistName)).Return(true, nil)
			},
			expectErr: errs.NewConflictError("Artist", "name", newArtistName),
		},
		{
			name: "Update_Error",
			req: dto.CreateArtistRequest{
				Name:  newArtistName,
				Image: &image,
			},
			prepareMock: func() {
				s.ArtistRepo.On("FindArtistById", mock.Anything, 1).Return(&models.Artist{
					Id:    1,
					Name:  oldArtistName,
					Slug:  utils.MakeSlug(oldArtistName),
					Image: utils.ParseImageToByte(&image),
				}, nil)

				s.ArtistRepo.On("FindExistsArtistBySlug", mock.Anything, utils.MakeSlug(newArtistName)).Return(false, nil)

				s.ArtistRepo.On("Update", mock.Anything, models.CreateArtistInput{
					Name:  newArtistName,
					Slug:  utils.MakeSlug(newArtistName),
					Image: utils.ParseImageToByte(&image),
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
			err := s.Svc.UpdateArtist(s.T().Context(), tc.req, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.ArtistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *ArtistServiceTestSuite) TestDeleteArtist() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.ArtistRepo.On("Delete", mock.Anything, 1).Return(nil)
			},
		},
		{
			name: "FindExistsArtistById_NotFound",
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Artist", "id", 1),
		},
		{
			name: "FindExistsArtistById_Error",
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(false, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "Delete_Error",
			prepareMock: func() {
				s.ArtistRepo.On("FindExistsArtistById", mock.Anything, 1).Return(true, nil)
				s.ArtistRepo.On("Delete", mock.Anything, 1).Return(errors.New("database failure"))
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
			err := s.Svc.DeleteArtist(s.T().Context(), 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.ArtistRepo.AssertExpectations(s.T())
		})
	}
}

func TestArtistServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ArtistServiceTestSuite))
}
