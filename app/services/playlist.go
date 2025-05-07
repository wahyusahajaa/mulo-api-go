package services

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type playlistService struct {
	repo contracts.PlaylistRepository
	log  *logrus.Logger
}

func NewPlaylistService(repo contracts.PlaylistRepository, log *logrus.Logger) contracts.PlaylistService {
	return &playlistService{
		repo: repo,
		log:  log,
	}
}

func (svc *playlistService) GetAll(ctx context.Context, pageSize int, offset int) (playlists []dto.Playlist, err error) {
	results, err := svc.repo.FindAll(ctx, pageSize, offset)
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

func (svc *playlistService) GetPlaylistById(ctx context.Context, id int) (playlist dto.Playlist, err error) {
	result, err := svc.repo.FindById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "GetPlaylistById", err)
		return playlist, err
	}
	if result == nil {
		notFoundErr := utils.NotFoundError{Resource: "Playlist", Id: id}
		utils.LogWarn(svc.log, ctx, "playlist_service", "GetPlaylistById", notFoundErr)
		return playlist, fmt.Errorf("%w", notFoundErr)
	}

	playlist.Id = result.Id
	playlist.Name = result.Name

	return playlist, nil
}

func (svc *playlistService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "GetCount", err)
		return
	}

	return
}

func (svc *playlistService) CreatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
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

func (svc *playlistService) UpdatePlaylist(ctx context.Context, req dto.CreatePlaylistRequest, id int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}

	exists, err := svc.repo.FindExistsPlaylistById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "UpdatePlaylist", err)
		return
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Playlist", Id: id}
		utils.LogWarn(svc.log, ctx, "playlist_service", "UpdatePlaylist", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	input := models.CreatePlaylistInput{
		Name:   req.Name,
		UserId: utils.GetUserId(ctx),
	}

	if err = svc.repo.Update(ctx, input, id); err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "UpdatePlaylist", err)
		return
	}

	return
}

func (svc *playlistService) DeletePlaylist(ctx context.Context, id int) (err error) {
	exists, err := svc.repo.FindExistsPlaylistById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "DeletePlaylist", err)
		return
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Playlist", Id: id}
		utils.LogWarn(svc.log, ctx, "playlist_service", "DeletePlaylist", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	if err = svc.repo.Delete(ctx, id); err != nil {
		utils.LogError(svc.log, ctx, "playlist_service", "DeletePlaylist", err)
		return
	}

	return
}
