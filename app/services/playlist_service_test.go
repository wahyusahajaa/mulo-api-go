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

const (
	userRole = "member"
	userId   = 1
	pageSize = 10
	offset   = 0
)

type PlaylistServiceTestSuite struct {
	suite.Suite
	Svc          contracts.PlaylistService
	playlistRepo *mocks.MockPlaylistRepository
	songRepo     *mocks.MockSongRepository
}

func (s *PlaylistServiceTestSuite) SetupTest() {
	s.playlistRepo = new(mocks.MockPlaylistRepository)
	s.songRepo = new(mocks.MockSongRepository)
	s.Svc = NewPlaylistService(s.playlistRepo, s.songRepo, nil)
}

func (s *PlaylistServiceTestSuite) ResetMocks() {
	s.playlistRepo.ExpectedCalls = nil
	s.playlistRepo.Calls = nil
	s.songRepo.ExpectedCalls = nil
	s.songRepo.Calls = nil
}

func (s *PlaylistServiceTestSuite) TestGetAll() {
	testCases := []struct {
		name          string
		prepareMock   func()
		expectResults []dto.Playlist
		expectTotal   int
		expectErr     error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.playlistRepo.On("FindCount", mock.Anything, userRole, userId).Return(1, nil)
				s.playlistRepo.On("FindAll", mock.Anything, userRole, userId, pageSize, offset).Return([]models.Playlist{
					{
						Id:   1,
						Name: "Playlist Test",
					},
				}, nil)
			},
			expectResults: []dto.Playlist{
				{
					Id:   1,
					Name: "Playlist Test",
				},
			},
			expectTotal: 1,
		},
		{
			name: "FindCount_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindCount", mock.Anything, userRole, userId).Return(0, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "FindAll_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindCount", mock.Anything, userRole, userId).Return(1, nil)
				s.playlistRepo.On("FindAll", mock.Anything, userRole, userId, pageSize, offset).Return(nil, errors.New("database failure"))
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
			results, total, err := s.Svc.GetAll(s.T().Context(), userRole, userId, pageSize, offset)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectResults, results)
				s.Equal(tc.expectTotal, total)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.playlistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PlaylistServiceTestSuite) TestGetPlaylistById() {
	testCases := []struct {
		name         string
		prepareMock  func()
		expectResult dto.Playlist
		expectErr    error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.playlistRepo.On("FindById", mock.Anything, userRole, userId, 1).Return(&models.Playlist{
					Id:   1,
					Name: "Playlist Test",
				}, nil)
			},
			expectResult: dto.Playlist{
				Id:   1,
				Name: "Playlist Test",
			},
		},
		{
			name: "FindById_NotFound",
			prepareMock: func() {
				s.playlistRepo.On("FindById", mock.Anything, userRole, userId, 1).Return(nil, nil)
			},
			expectErr: errs.NewNotFoundError("Playlist", "id", 1),
		},
		{
			name: "FindById_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindById", mock.Anything, userRole, userId, 1).Return(nil, errors.New("database failure"))
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
			result, err := s.Svc.GetPlaylistById(s.T().Context(), userRole, userId, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectResult, result)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.playlistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PlaylistServiceTestSuite) TestCreatePlaylist() {
	var validationErr = errors.New("validation failed")

	testCases := []struct {
		name            string
		req             dto.CreatePlaylistRequest
		prepareMock     func()
		expectErr       error
		expectValErrMap map[string]string
	}{
		{
			name: "success",
			req: dto.CreatePlaylistRequest{
				Name: "Test Playlist",
			},
			prepareMock: func() {
				s.playlistRepo.On("Store", mock.Anything, models.CreatePlaylistInput{
					Name: "Test Playlist",
				}).Return(nil)
			},
		},
		{
			name: "ValidationErrors_RequiredName",
			req: dto.CreatePlaylistRequest{
				Name: "",
			},
			expectErr:       validationErr,
			expectValErrMap: map[string]string{"name": "Field is required"},
		},
		{
			name: "Store_Error",
			req: dto.CreatePlaylistRequest{
				Name: "Test Playlist",
			},
			prepareMock: func() {
				s.playlistRepo.On("Store", mock.Anything, models.CreatePlaylistInput{
					Name: "Test Playlist",
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
			err := s.Svc.CreatePlaylist(s.T().Context(), tc.req)

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

			s.playlistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PlaylistServiceTestSuite) TestUpdatePlaylist() {
	validationErr := errors.New("validation failed")

	testCases := []struct {
		name            string
		req             dto.CreatePlaylistRequest
		prepareMock     func()
		expectErr       error
		expectValErrMap map[string]string
	}{
		{
			name: "success",
			req: dto.CreatePlaylistRequest{
				Name: "Test Playlist",
			},
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.playlistRepo.On("Update", mock.Anything, models.CreatePlaylistInput{
					Name: "Test Playlist",
				}, 1).Return(nil)
			},
		},
		{
			name: "ValidationErrors_RequiredName",
			req: dto.CreatePlaylistRequest{
				Name: "",
			},
			expectErr:       validationErr,
			expectValErrMap: map[string]string{"name": "Field is required"},
		},
		{
			name: "FindExistsPlaylistById_NotFound",
			req: dto.CreatePlaylistRequest{
				Name: "Test Playlist",
			},
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Playlist", "id", 1),
		},
		{
			name: "FindExistsPlaylistById_Error",
			req: dto.CreatePlaylistRequest{
				Name: "Test Playlist",
			},
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(false, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "Update_Error",
			req: dto.CreatePlaylistRequest{
				Name: "Test Playlist",
			},
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.playlistRepo.On("Update", mock.Anything, models.CreatePlaylistInput{
					Name: "Test Playlist",
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
			err := s.Svc.UpdatePlaylist(s.T().Context(), tc.req, userRole, userId, 1)

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

			s.playlistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PlaylistServiceTestSuite) TestDeletePlaylist() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.playlistRepo.On("Delete", mock.Anything, userRole, userId, 1).Return(nil)
			},
		},
		{
			name: "FindExistsPlaylistById_NotFound",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(nil, nil)
			},
			expectErr: errs.NewNotFoundError("Playlist", "id", 1),
		},
		{
			name: "FindExistsPlaylistById_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(nil, errors.New("database failure"))
			},
			expectErr: errors.New("database failure"),
		},
		{
			name: "Delete_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.playlistRepo.On("Delete", mock.Anything, userRole, userId, 1).Return(errors.New("database failure"))
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
			err := s.Svc.DeletePlaylist(s.T().Context(), userRole, userId, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.playlistRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PlaylistServiceTestSuite) TestGetPlaylistSongs() {
	image := dto.Image{Src: "image.png", BlurHash: "abc"}
	imageBytes := utils.ParseImageToByte(&image)

	testCases := []struct {
		name          string
		prepareMock   func()
		expectResults []dto.Song
		expectErr     error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.playlistRepo.On("FindById", mock.Anything, userRole, userId, 1).Return(&models.Playlist{
					Id:   1,
					Name: "Playlist Test",
				}, nil)
				s.playlistRepo.On("FindPlaylistSongs", mock.Anything, 1, pageSize, offset).Return([]models.Song{
					{
						Id:       1,
						AlbumId:  1,
						Title:    "Aku pulang",
						Audio:    "akupulang.mp3",
						Duration: 352,
						Image:    imageBytes,
						Album: models.AlbumWithArtist{
							Album: models.Album{
								Id:       1,
								ArtistId: 1,
								Name:     "Album test",
								Slug:     "album-test",
								Image:    imageBytes,
							},
							Artist: models.Artist{
								Id:    1,
								Name:  "Sheila",
								Slug:  "sheila",
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
					Audio:    "akupulang.mp3",
					Duration: 352,
					Image:    image,
					Album: dto.AlbumWithArtist{
						Album: dto.Album{
							Id:    1,
							Name:  "Album test",
							Slug:  "album-test",
							Image: image,
						},
						Artist: dto.Artist{
							Id:    1,
							Name:  "Sheila",
							Slug:  "sheila",
							Image: image,
						},
					},
				},
			},
		},
		{
			name: "FindById_NotFound",
			prepareMock: func() {
				s.playlistRepo.On("FindById", mock.Anything, userRole, userId, 1).Return(nil, nil)
			},
			expectErr: errs.NewNotFoundError("Playlist", "id", 1),
		},
		{
			name: "FindPlaylistSongs_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindById", mock.Anything, userRole, userId, 1).Return(&models.Playlist{
					Id:   1,
					Name: "Playlist Test",
				}, nil)
				s.playlistRepo.On("FindPlaylistSongs", mock.Anything, 1, pageSize, offset).Return(nil, errors.New("database failure"))
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
			results, err := s.Svc.GetPlaylistSongs(s.T().Context(), userRole, userId, 1, pageSize, offset)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
				s.Equal(tc.expectResults, results)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.playlistRepo.AssertExpectations(s.T())
			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PlaylistServiceTestSuite) TestCreatePlaylistSong() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.playlistRepo.On("FindExistsPlaylistSong", mock.Anything, 1, 1).Return(false, nil)
				s.playlistRepo.On("StorePlaylistSong", mock.Anything, 1, 1).Return(nil)
			},
		},
		{
			name: "FindExistsPlaylistById_NotFound",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Playlist", "id", 1),
		},
		{
			name: "FindExistsSongById_NotFound",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Song", "id", 1),
		},
		{
			name: "FindExistsPlaylistSong_Conflict",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.playlistRepo.On("FindExistsPlaylistSong", mock.Anything, 1, 1).Return(true, nil)
			},
			expectErr: errs.NewConflictError("PlaylistSong", "song_id", 1),
		},
		{
			name: "StorePlaylistSong_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.songRepo.On("FindExistsSongById", mock.Anything, 1).Return(true, nil)
				s.playlistRepo.On("FindExistsPlaylistSong", mock.Anything, 1, 1).Return(false, nil)
				s.playlistRepo.On("StorePlaylistSong", mock.Anything, 1, 1).Return(errors.New("database failure"))
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
			err := s.Svc.CreatePlaylistSong(s.T().Context(), userRole, userId, 1, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.playlistRepo.AssertExpectations(s.T())
			s.songRepo.AssertExpectations(s.T())
		})
	}
}

func (s *PlaylistServiceTestSuite) TestDeletePlaylistSong() {
	testCases := []struct {
		name        string
		prepareMock func()
		expectErr   error
	}{
		{
			name: "success",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.playlistRepo.On("FindExistsPlaylistSong", mock.Anything, 1, 1).Return(true, nil)
				s.playlistRepo.On("DeletePlaylistSong", mock.Anything, 1, 1).Return(nil)
			},
		},
		{
			name: "FindExistsPlaylistById_NotFound",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("Playlist", "id", 1),
		},
		{
			name: "FindExistsPlaylistSong_NotFound",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.playlistRepo.On("FindExistsPlaylistSong", mock.Anything, 1, 1).Return(false, nil)
			},
			expectErr: errs.NewNotFoundError("PlaylistSong", "song_id", 1),
		},
		{
			name: "DeletePlaylistSong_Error",
			prepareMock: func() {
				s.playlistRepo.On("FindExistsPlaylistById", mock.Anything, userRole, userId, 1).Return(true, nil)
				s.playlistRepo.On("FindExistsPlaylistSong", mock.Anything, 1, 1).Return(true, nil)
				s.playlistRepo.On("DeletePlaylistSong", mock.Anything, 1, 1).Return(errors.New("database failure"))
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
			err := s.Svc.DeletePlaylistSong(s.T().Context(), userRole, userId, 1, 1)

			// Assert
			if tc.expectErr == nil {
				s.NoError(err)
			} else {
				s.Error(err)
				s.EqualError(err, tc.expectErr.Error())
			}

			s.playlistRepo.AssertExpectations(s.T())
		})
	}
}

func TestPlaylistServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PlaylistServiceTestSuite))
}
