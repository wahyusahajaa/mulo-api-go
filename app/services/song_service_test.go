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

type SongServiceTestSuite struct {
	suite.Suite
	Svc       contracts.SongService
	songRepo  *mocks.MockSongRepository
	albumRepo *mocks.MockAlbumRepository
}

func (s *SongServiceTestSuite) SetupTest() {
	s.songRepo = new(mocks.MockSongRepository)
	s.albumRepo = new(mocks.MockAlbumRepository)
	s.Svc = NewSongService(s.songRepo, s.albumRepo, nil)
}

func (s *SongServiceTestSuite) ResetMocks() {
	s.songRepo.ExpectedCalls = nil
	s.songRepo.Calls = nil
	s.albumRepo.ExpectedCalls = nil
	s.albumRepo.Calls = nil
}

func (s *SongServiceTestSuite) TestGetAll() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abc"}
	image1Bytes := utils.ParseImageToByte(&image1)

	testCases := []struct {
		name            string
		prepareMock     func()
		expectedResults []dto.Song
		expectedTotal   int
		expectedErr     error
	}{
		{
			name: "success",
			prepareMock: func() {
				// Setup mocks
				s.songRepo.On("FindCount", mock.Anything).Return(1, nil)
				s.songRepo.On("FindAll", mock.Anything, 10, 0).Return([]models.Song{
					{
						Id:       1,
						AlbumId:  1,
						Title:    "Aku pulang",
						Audio:    "aku-pulang.mp3",
						Duration: 352,
						Image:    image1Bytes,
						Album: models.AlbumWithArtist{
							Album: models.Album{
								Id:       1,
								ArtistId: 1,
								Name:     "Album naff",
								Slug:     "album-naff",
								Image:    image1Bytes,
							},
							Artist: models.Artist{
								Id:    1,
								Name:  "Naff",
								Slug:  "naff",
								Image: image1Bytes,
							},
						},
					},
				}, nil)
			},
			expectedResults: []dto.Song{
				{
					Id:       1,
					Title:    "Aku pulang",
					Audio:    "aku-pulang.mp3",
					Duration: 352,
					Image:    image1,
					Album: dto.AlbumWithArtist{
						Album: dto.Album{
							Id:    1,
							Name:  "Album naff",
							Slug:  "album-naff",
							Image: image1,
						},
						Artist: dto.Artist{
							Id:    1,
							Name:  "Naff",
							Slug:  "naff",
							Image: image1,
						},
					},
				},
			},
			expectedTotal: 1,
		},
		{
			name: "FindCountError",
			prepareMock: func() {
				s.songRepo.On("FindCount", mock.Anything).Return(0, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "findAllError",
			prepareMock: func() {
				s.songRepo.On("FindCount", mock.Anything).Return(1, nil)
				s.songRepo.On("FindAll", mock.Anything, 10, 0).Return(nil, errors.New("database failed"))
			},
			expectedErr: errors.New("database failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			// Actual
			results, total, err := s.Svc.GetAll(s.T().Context(), 10, 0)

			// Assert
			if tc.expectedErr == nil {
				s.NoError(err)
				s.Equal(tc.expectedTotal, total)
				s.Equal(tc.expectedResults, results)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func (s *SongServiceTestSuite) TestGetSongById() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abc"}
	image1Bytes := utils.ParseImageToByte(&image1)

	testCases := []struct {
		name           string
		prepareMock    func()
		expectedResult dto.Song
		expectedErr    error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.songRepo.On("FindSongById", mock.Anything, 1).Return(&models.Song{
					Id:       1,
					AlbumId:  1,
					Title:    "Aku pulang",
					Audio:    "aku-pulang.mp3",
					Duration: 352,
					Image:    image1Bytes,
					Album: models.AlbumWithArtist{
						Album: models.Album{
							Id:       1,
							ArtistId: 1,
							Name:     "Album naff",
							Slug:     "album-naff",
							Image:    image1Bytes,
						},
						Artist: models.Artist{
							Id:    1,
							Name:  "Naff",
							Slug:  "naff",
							Image: image1Bytes,
						},
					},
				}, nil)
			},
			expectedResult: dto.Song{
				Id:       1,
				Title:    "Aku pulang",
				Audio:    "aku-pulang.mp3",
				Duration: 352,
				Image:    image1,
				Album: dto.AlbumWithArtist{
					Album: dto.Album{
						Id:    1,
						Name:  "Album naff",
						Slug:  "album-naff",
						Image: image1,
					},
					Artist: dto.Artist{
						Id:    1,
						Name:  "Naff",
						Slug:  "naff",
						Image: image1,
					},
				},
			},
		},
		{
			name: "FindSongById_NotFound",
			prepareMock: func() {
				s.songRepo.On("FindSongById", mock.Anything, 1).Return(nil, nil)
			},
			expectedErr: errs.NewNotFoundError("Song", "id", 1),
		},
		{
			name: "FindSongById_Error",
			prepareMock: func() {
				s.songRepo.On("FindSongById", mock.Anything, 1).Return(nil, errors.New("database failed"))
			},
			expectedErr: errors.New("database failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			// Actual
			result, err := s.Svc.GetSongById(s.T().Context(), 1)

			// Assert
			if tc.expectedErr == nil {
				s.NoError(err)
				s.Equal(tc.expectedResult, result)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func (s *SongServiceTestSuite) TestCreateSong() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abc"}
	var validationErr = errors.New("validation failed")

	testCases := []struct {
		name              string
		createSongRequest dto.CreateSongRequest
		prepareMock       func()
		expectedErr       error
		expectedValErrMap map[string]string
	}{
		{
			name: "success",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image1,
			},
			prepareMock: func() {
				s.albumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(true, nil)
				s.songRepo.On("Store", mock.Anything, models.CreateSongInput{
					AlbumId:  1,
					Title:    "Aku bukanlah superman",
					Audio:    "song-1.mp3",
					Duration: 350,
					Image:    utils.ParseImageToByte(&image1),
				}).Return(nil)
			},
		},
		{
			name: "ValidationFailed_AllRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  0,
				Title:    "",
				Audio:    "",
				Duration: 0,
			},
			expectedErr: validationErr,
			expectedValErrMap: map[string]string{
				"album_id": "Field is required",
				"title":    "Field is required",
				"audio":    "Field is required",
				"duration": "Field is required",
				"image":    "Field is required",
			},
		},
		{
			name: "ValidationFailed_AlbumIdRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  0,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image1,
			},
			expectedErr:       validationErr,
			expectedValErrMap: map[string]string{"album_id": "Field is required"},
		},
		{
			name: "ValidationFailed_TitleRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image1,
			},
			expectedErr:       validationErr,
			expectedValErrMap: map[string]string{"title": "Field is required"},
		},
		{
			name: "ValidationFailed_AudioRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "",
				Duration: 350,
				Image:    &image1,
			},
			expectedErr:       validationErr,
			expectedValErrMap: map[string]string{"audio": "Field is required"},
		},
		{
			name: "ValidationFailed_DurationRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 0,
				Image:    &image1,
			},
			expectedErr:       validationErr,
			expectedValErrMap: map[string]string{"duration": "Field is required"},
		},
		{
			name: "FindExistsAlbumById_Error",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image1,
			},
			prepareMock: func() {
				s.albumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(false, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "Store_Error",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image1,
			},
			prepareMock: func() {
				s.albumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(true, nil)
				s.songRepo.On("Store", mock.Anything, models.CreateSongInput{
					AlbumId:  1,
					Title:    "Aku bukanlah superman",
					Audio:    "song-1.mp3",
					Duration: 350,
					Image:    utils.ParseImageToByte(&image1),
				}).Return(errors.New("database failure"))
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

			// Actual
			err := s.Svc.CreateSong(s.T().Context(), tc.createSongRequest)

			if tc.expectedErr == nil {
				s.NoError(err)
			} else if errors.Is(tc.expectedErr, validationErr) {
				var valErr *errs.BadRequestError
				s.ErrorAs(err, &valErr)
				s.Equal(tc.expectedValErrMap, valErr.Errors)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.albumRepo.AssertExpectations(s.T())
			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func (s *SongServiceTestSuite) TestUpdateSong() {
	image := dto.Image{Src: "image1.png", BlurHash: "abc"}
	var validationErr = errors.New("validation failed")

	testCases := []struct {
		name                string
		createSongRequest   dto.CreateSongRequest
		prepareMock         func()
		expectedErr         error
		expectedValErrorMap map[string]string
	}{
		{
			name: "success",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.albumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(true, nil)
				s.songRepo.On("Update", mock.Anything, models.CreateSongInput{
					AlbumId:  1,
					Title:    "Aku bukanlah superman",
					Audio:    "song-1.mp3",
					Duration: 350,
					Image:    utils.ParseImageToByte(&image),
				}, 1).Return(nil)
			},
		},
		{
			name: "ValidationFailed_AllRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  0,
				Title:    "",
				Audio:    "",
				Duration: 0,
			},
			expectedErr: validationErr,
			expectedValErrorMap: map[string]string{
				"album_id": "Field is required",
				"title":    "Field is required",
				"audio":    "Field is required",
				"duration": "Field is required",
				"image":    "Field is required",
			},
		},
		{
			name: "ValidationFailed_AlbumIdRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  0,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			expectedErr:         validationErr,
			expectedValErrorMap: map[string]string{"album_id": "Field is required"},
		},
		{
			name: "ValidationFailed_TitleRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			expectedErr:         validationErr,
			expectedValErrorMap: map[string]string{"title": "Field is required"},
		},
		{
			name: "ValidationFailed_AudioRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "",
				Duration: 350,
				Image:    &image,
			},
			expectedErr:         validationErr,
			expectedValErrorMap: map[string]string{"audio": "Field is required"},
		},
		{
			name: "ValidationFailed_DurationRequired",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 0,
				Image:    &image,
			},
			expectedErr:         validationErr,
			expectedValErrorMap: map[string]string{"duration": "Field is required"},
		},
		{
			name: "FindExistsSongById_NotFound",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(false, nil)
			},
			expectedErr: errs.NewNotFoundError("Song", "id", 1),
		},
		{
			name: "FindExistsSongById_Error",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(false, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "FindExistsAlbumById_NotFound",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.albumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(false, nil)
			},
			expectedErr: errs.NewNotFoundError("Album", "id", 1),
		},
		{
			name: "FindExistsAlbumById_Error",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.albumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(false, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "Update_Error",
			createSongRequest: dto.CreateSongRequest{
				AlbumId:  1,
				Title:    "Aku bukanlah superman",
				Audio:    "song-1.mp3",
				Duration: 350,
				Image:    &image,
			},
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.albumRepo.On("FindExistsAlbumById", mock.Anything, 1).Return(true, nil)
				s.songRepo.On("Update", mock.Anything, models.CreateSongInput{
					AlbumId:  1,
					Title:    "Aku bukanlah superman",
					Audio:    "song-1.mp3",
					Duration: 350,
					Image:    utils.ParseImageToByte(&image),
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

			// Actual
			err := s.Svc.UpdateSong(s.T().Context(), tc.createSongRequest, 1)

			// Assert
			if tc.expectedErr == nil {
				s.NoError(err)
			} else if errors.Is(tc.expectedErr, validationErr) {
				var valErr *errs.BadRequestError
				s.ErrorAs(err, &valErr)
				s.Equal(tc.expectedValErrorMap, valErr.Errors)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.songRepo.AssertExpectations(s.T())
			s.albumRepo.AssertExpectations(s.T())
		})
	}

}

func (s *SongServiceTestSuite) TestDeleteSong() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectedErr error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.songRepo.On("Delete", mock.Anything, 1).Return(nil)
			},
		},
		{
			name: "FindExistsSongById_NotFound",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(false, nil)
			},
			expectedErr: errs.NewNotFoundError("Song", "id", 1),
		},
		{
			name: "FindExistsSongById_Error",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(false, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "Delete_Error",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.songRepo.On("Delete", mock.Anything, 1).Return(errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			// Actual
			err := s.Svc.DeleteSong(s.T().Context(), 1)

			// Assert
			if tc.expectedErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func (s *SongServiceTestSuite) TestGetSongsByAlbumId() {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abc"}
	image1Bytes := utils.ParseImageToByte(&image1)

	testCases := []struct {
		name            string
		prepareMock     func()
		expectedResults []dto.Song
		expectedTotal   int
		expectedErr     error
	}{
		{
			name: "success",
			prepareMock: func() {
				// Setup mocks
				s.songRepo.On("FindCountSongsByAlbumId", mock.Anything, 1).Return(1, nil)
				s.songRepo.On("FindSongsByAlbumId", mock.Anything, 1, 10, 0).Return([]models.Song{
					{
						Id:       1,
						AlbumId:  1,
						Title:    "Aku pulang",
						Audio:    "aku-pulang.mp3",
						Duration: 352,
						Image:    image1Bytes,
						Album: models.AlbumWithArtist{
							Album: models.Album{
								Id:       1,
								ArtistId: 1,
								Name:     "Album naff",
								Slug:     "album-naff",
								Image:    image1Bytes,
							},
							Artist: models.Artist{
								Id:    1,
								Name:  "Naff",
								Slug:  "naff",
								Image: image1Bytes,
							},
						},
					},
				}, nil)
			},
			expectedResults: []dto.Song{
				{
					Id:       1,
					Title:    "Aku pulang",
					Audio:    "aku-pulang.mp3",
					Duration: 352,
					Image:    image1,
					Album: dto.AlbumWithArtist{
						Album: dto.Album{
							Id:    1,
							Name:  "Album naff",
							Slug:  "album-naff",
							Image: image1,
						},
						Artist: dto.Artist{
							Id:    1,
							Name:  "Naff",
							Slug:  "naff",
							Image: image1,
						},
					},
				},
			},
			expectedTotal: 1,
		},
		{
			name: "FindCountSongsByAlbumId_Error",
			prepareMock: func() {
				s.songRepo.On("FindCountSongsByAlbumId", mock.Anything, 1).Return(0, errors.New("database failure"))
			},
			expectedErr: errors.New("database failure"),
		},
		{
			name: "FindSongsByAlbumId_Error",
			prepareMock: func() {
				s.songRepo.On("FindCountSongsByAlbumId", mock.Anything, 1).Return(1, nil)
				s.songRepo.On("FindSongsByAlbumId", mock.Anything, 1, 10, 0).Return(nil, errors.New("database failed"))
			},
			expectedErr: errors.New("database failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ResetMocks()
			tc.prepareMock()

			// Actual
			results, total, err := s.Svc.GetSongsByAlbumId(s.T().Context(), 1, 10, 0)

			// Assert
			if tc.expectedErr == nil {
				s.NoError(err)
				s.Equal(tc.expectedTotal, total)
				s.Equal(tc.expectedResults, results)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectedErr.Error())
			}

			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func TestSongServiceTestSuite(t *testing.T) {
	suite.Run(t, new(SongServiceTestSuite))
}
