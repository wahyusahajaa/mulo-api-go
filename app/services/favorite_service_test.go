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

type FavoriteServiceTestSuite struct {
	suite.Suite
	Svc      contracts.FavoriteService
	favRepo  *mocks.MockFavoriteRepository
	songRepo *mocks.MockSongRepository
}

func (s *FavoriteServiceTestSuite) SetupTest() {
	s.favRepo = new(mocks.MockFavoriteRepository)
	s.songRepo = new(mocks.MockSongRepository)
	s.Svc = NewFavoriteService(s.favRepo, s.songRepo, nil)
}

func (s *FavoriteServiceTestSuite) ResetMocks() {
	s.favRepo.Calls = nil
	s.favRepo.ExpectedCalls = nil
	s.songRepo.Calls = nil
	s.songRepo.ExpectedCalls = nil
}

func (s *FavoriteServiceTestSuite) TestGetFavoriteSongsByUserID() {
	image := dto.Image{Src: "image.png", BlurHash: "abc"}
	imageBytes := utils.ParseImageToByte(&image)

	testCases := []struct {
		name          string
		prepareMock   func()
		expectResults []dto.Song
		expectTotal   int
		expectErr     error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.favRepo.On("FindCountFavoriteSongsByUserID", mock.Anything, 1).Return(1, nil)
				s.favRepo.On("FindFavoriteSongsByUserID", mock.Anything, 1, 10, 0).Return([]models.Song{
					{
						Id:       1,
						AlbumId:  1,
						Title:    "Aku pulang",
						Audio:    "aku.mp3",
						Duration: 352,
						Image:    imageBytes,
						Album: models.AlbumWithArtist{
							Album: models.Album{
								Id:       1,
								ArtistId: 1,
								Name:     "Album aku pulang",
								Slug:     "album-aku-pulang",
								Image:    imageBytes,
							},
							Artist: models.Artist{
								Id:    1,
								Name:  "Noah",
								Slug:  "noah",
								Image: imageBytes,
							},
						},
					},
				}, nil)
			},
			expectResults: []dto.Song{
				{
					Id:       1,
					Title:    "Aku pulang",
					Audio:    "aku.mp3",
					Duration: 352,
					Image:    image,
					Album: dto.AlbumWithArtist{
						Album: dto.Album{
							Id:    1,
							Name:  "Album aku pulang",
							Slug:  "album-aku-pulang",
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
			},
			expectTotal: 1,
		},
		{
			name: "FindCountFavoriteSongsByUserID_Error",
			prepareMock: func() {
				s.favRepo.On("FindCountFavoriteSongsByUserID", mock.Anything, 1).Return(0, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "FindFavoriteSongsByUserID_Error",
			prepareMock: func() {
				s.favRepo.On("FindCountFavoriteSongsByUserID", mock.Anything, 1).Return(1, nil)
				s.favRepo.On("FindFavoriteSongsByUserID", mock.Anything, 1, 10, 0).Return(nil, errors.New("database failure"))
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
			results, total, err := s.Svc.GetFavoriteSongsByUserID(s.T().Context(), 1, 10, 0)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectResults, results)
				s.Equal(tc.expectTotal, total)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.favRepo.AssertExpectations(s.T())
		})
	}
}

func (s *FavoriteServiceTestSuite) TestAddFavoriteSong() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.favRepo.On("FindExistsFavoriteSongBySongID", mock.Anything, 1, 1).Return(false, nil)
				s.favRepo.On("StoreFavoriteSong", mock.Anything, 1, 1).Return(nil)
			},
		},
		{
			name: "FindExistsSongById_NotFound",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Song", "id", 1),
		},
		{
			name: "FindExistsFavoriteSongBySongID_Conflict",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.favRepo.On("FindExistsFavoriteSongBySongID", mock.Anything, 1, 1).Return(true, nil)
			},
			expectErr: errs.NewConflictError("Song on favorite", "song_id", 1),
		},
		{
			name: "StoreFavoriteSong_Error",
			prepareMock: func() {
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.favRepo.On("FindExistsFavoriteSongBySongID", mock.Anything, 1, 1).Return(false, nil)
				s.favRepo.On("StoreFavoriteSong", mock.Anything, 1, 1).Return(errors.New("database failure"))
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
			err := s.Svc.AddFavoriteSong(s.T().Context(), 1, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.favRepo.AssertExpectations(s.T())
			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func (s *FavoriteServiceTestSuite) TestRemoveFavoriteSong() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.favRepo.On("FindExistsFavoriteSongBySongID", mock.Anything, 1, 1).Return(true, nil)
				s.favRepo.On("DeleteFavoriteSong", mock.Anything, 1, 1).Return(nil)
			},
		},
		{
			name: "FindExistsFavoriteSongBySongID_NotFound",
			prepareMock: func() {
				s.favRepo.On("FindExistsFavoriteSongBySongID", mock.Anything, 1, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Song on favorite", "song_id", 1),
		},
		{
			name: "DeleteFavoriteSong_Error",
			prepareMock: func() {
				s.favRepo.On("FindExistsFavoriteSongBySongID", mock.Anything, 1, 1).Return(true, nil)
				s.favRepo.On("DeleteFavoriteSong", mock.Anything, 1, 1).Return(errors.New("database failure"))
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
			err := s.Svc.RemoveFavoriteSong(s.T().Context(), 1, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.favRepo.AssertExpectations(s.T())
		})
	}
}

func TestFavoriteServiceTestSuite(t *testing.T) {
	suite.Run(t, new(FavoriteServiceTestSuite))
}
