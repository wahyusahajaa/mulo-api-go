package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type playlistService struct {
	repo     contracts.PlaylistRepository
	songRepo contracts.SongRepository
	log      *logrus.Logger
}

func NewPlaylistService(repo contracts.PlaylistRepository, songRepo contracts.SongRepository, log *logrus.Logger) contracts.PlaylistService {
	return &playlistService{
		repo:     repo,
		songRepo: songRepo,
		log:      log,
	}
}

func (svc *playlistService) GetAll(ctx context.Context, userRole string, userId, pageSize, offset int) (playlists []dto.Playlist, err error) {
	results, err := svc.repo.FindAll(ctx, userRole, userId, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "GetAll", err)
		return nil, err
	}

	playlists = make([]dto.Playlist, 0, len(results))
	for _, result := range results {
		playlists = append(playlists, dto.Playlist{
			Id:   result.Id,
			Name: result.Name,
		})
	}

	return playlists, nil
}

func (svc *playlistService) GetPlaylistById(ctx context.Context, userRole string, userId, id int) (playlist dto.Playlist, err error) {
	result, err := svc.repo.FindById(ctx, userRole, userId, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "GetPlaylistById", err)
		return playlist, err
	}
	if result == nil {
		notFoundErr := errs.NewNotFoundError("Playlist", "id", id)
		utils.LogWarn(svc.log, ctx, "playlist_service", "GetPlaylistById", notFoundErr)
		return playlist, notFoundErr
	}

	playlist.Id = result.Id
	playlist.Name = result.Name

	return playlist, nil
}

func (svc *playlistService) GetCount(ctx context.Context, userRole string, userId int) (total int, err error) {
	total, err = svc.repo.FindCount(ctx, userRole, userId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "GetCount", err)
		return
	}

	return
}

func (svc *playlistService) CreatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	input := models.CreatePlaylistInput{
		UserId: utils.GetUserId(ctx),
		Name:   req.Name,
	}

	if err = svc.repo.Store(ctx, input); err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "CreatePlaylist", err)
		return
	}

	return
}

func (svc *playlistService) UpdatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest, userRole string, userId, playlistId int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	exists, err := svc.repo.FindExistsPlaylistById(ctx, userRole, userId, playlistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "UpdatePlaylist", err)
		return
	}

	if !exists {
		notFoundErr := errs.NewNotFoundError("Playlist", "id", playlistId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "UpdatePlaylist", notFoundErr)
		return notFoundErr
	}

	input := models.CreatePlaylistInput{
		Name:   req.Name,
		UserId: utils.GetUserId(ctx),
	}

	if err = svc.repo.Update(ctx, input, playlistId); err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "UpdatePlaylist", err)
		return
	}

	return
}

func (svc *playlistService) DeletePlaylist(ctx context.Context, userRole string, userId, playlistId int) (err error) {
	exists, err := svc.repo.FindExistsPlaylistById(ctx, userRole, userId, playlistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "DeletePlaylist", err)
		return
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Playlist", "id", playlistId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "DeletePlaylist", notFoundErr)
		return notFoundErr
	}

	if err = svc.repo.Delete(ctx, userRole, userId, playlistId); err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "DeletePlaylist", err)
		return
	}

	return
}

func (svc *playlistService) GetPlaylistSongs(ctx context.Context, role string, userId, playlistId, pageSize, offset int) (songs []dto.Song, err error) {
	playlist, err := svc.repo.FindById(ctx, role, userId, playlistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "GetPlaylistSongs", err)
		return nil, err
	}
	if playlist == nil {
		notfoundErr := errs.NewNotFoundError("Playlist", "id", playlistId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "GetPlaylistSongs", notfoundErr)
		return nil, notfoundErr
	}

	results, err := svc.repo.FindPlaylistSongs(ctx, playlist.Id, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "GetPlaylistSongs", err)
		return nil, err
	}

	songs = make([]dto.Song, 0, len(results))
	for _, v := range results {
		song := dto.Song{
			Id:       v.Id,
			Title:    v.Title,
			Audio:    v.Audio,
			Duration: v.Duration,
			Image:    utils.ParseImageToJSON(v.Image),
			Album: dto.AlbumWithArtist{
				Album: dto.Album{
					Id:    v.Album.Id,
					Name:  v.Album.Name,
					Slug:  v.Album.Slug,
					Image: utils.ParseImageToJSON(v.Album.Image),
				},
				Artist: dto.Artist{
					Id:    v.Album.Artist.Id,
					Name:  v.Album.Artist.Name,
					Slug:  v.Album.Artist.Slug,
					Image: utils.ParseImageToJSON(v.Album.Artist.Image),
				},
			},
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func (svc *playlistService) CreatePlaylistSong(ctx context.Context, userRole string, userId, playlistId, songId int) (err error) {
	// Check existing playlist
	exists, err := svc.repo.FindExistsPlaylistById(ctx, userRole, userId, playlistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "CreatePlaylistSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewConflictError("Playlist", "id", playlistId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "CreatePlaylistSong", notFoundErr)
		return notFoundErr
	}

	// Check existing song
	exists, err = svc.songRepo.FindExistsSongById(ctx, songId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "CreatePlaylistSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Song", "id", songId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "CreatePlaylistSong", notFoundErr)
		return notFoundErr
	}

	// Check if playlist song already exists with the same song id
	exists, err = svc.repo.FindExistsPlaylistSong(ctx, playlistId, songId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "CreatePlaylistSong", err)
		return err
	}
	if exists {
		conflictErr := errs.NewConflictError("PlaylistSong", "song_id", songId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "CreatePlaylistSong", conflictErr)
		return conflictErr
	}

	if err := svc.repo.StorePlaylistSong(ctx, playlistId, songId); err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "CreatePlaylistSong", err)
		return err
	}

	return
}

func (svc *playlistService) DeletePlaylistSong(ctx context.Context, userRole string, userId, playlistId, songId int) (err error) {
	// Check existing playlist
	exists, err := svc.repo.FindExistsPlaylistById(ctx, userRole, userId, playlistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "DeletePlaylistSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Playlist", "id", playlistId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "DeletePlaylistSong", notFoundErr)
		return notFoundErr
	}

	// Check existing song on playlists
	exists, err = svc.repo.FindExistsPlaylistSong(ctx, playlistId, songId)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "DeletePlaylistSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("PlaylistSong", "song_id", songId)
		utils.LogWarn(svc.log, ctx, "playlist_service", "DeletePlaylistSong", notFoundErr)
		return notFoundErr
	}

	if err := svc.repo.DeletePlaylistSong(ctx, playlistId, songId); err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "DeletePlaylistSong", err)
		return err
	}

	return
}
