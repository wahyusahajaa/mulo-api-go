package services

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type favoriteService struct {
	favRepo  contracts.FavoriteRepository
	songRepo contracts.SongRepository
	log      *logrus.Logger
}

func NewFavoriteService(favRepo contracts.FavoriteRepository, songRepo contracts.SongRepository, log *logrus.Logger) contracts.FavoriteService {
	return &favoriteService{
		favRepo:  favRepo,
		songRepo: songRepo,
		log:      log,
	}
}

func (svc *favoriteService) CreateFavoriteSong(ctx context.Context, userId int, songId int) (err error) {
	exists, err := svc.songRepo.FindExistsSongById(ctx, songId)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "CreateFavorite", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Song", Id: songId}
		utils.LogWarn(svc.log, ctx, "favorite_service", "CreateFavorite", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	exists, err = svc.favRepo.FindExistsFavoriteSongBySongId(ctx, userId, songId)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "CreateFavorite", err)
		return err
	}

	if exists {
		conflictErr := utils.ConflictError{Resource: "Favorite", Field: "song_id", Value: songId}
		utils.LogWarn(svc.log, ctx, "favorite_service", "CreateFavorite", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	if err = svc.favRepo.StoreFavoriteSong(ctx, userId, songId); err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "CreateFavorite", err)
		return err
	}

	return
}

func (svc *favoriteService) GetCountFavoriteSongsByUserId(ctx context.Context, userId int) (total int, err error) {
	total, err = svc.favRepo.FindCountFavoriteSongsByUserId(ctx, userId)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "GetCountFavoriteSongsByUserId", err)
		return
	}

	return
}

func (svc *favoriteService) GetAllFavoriteSongsByUserId(ctx context.Context, userId, pageSize, offset int) (songs []dto.Song, err error) {
	results, err := svc.favRepo.FindFavoriteSongsByUserId(ctx, userId, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "GetAllFavoriteSongsByUserId", err)
		return nil, err
	}

	songs = make([]dto.Song, 0, len(results))
	for _, result := range results {
		var song dto.Song
		song.Id = result.Id
		song.Title = result.Title
		song.Duration = result.Duration
		song.Audio = result.Audio
		song.Image = utils.ParseImageToJSON(result.Image)
		song.Album = dto.Album{
			Id:    result.Album.Id,
			Name:  result.Album.Name,
			Slug:  result.Album.Slug,
			Image: utils.ParseImageToJSON(result.Album.Image),
			Artist: dto.Artist{
				Id:    result.Album.Artist.Id,
				Name:  result.Album.Artist.Name,
				Slug:  result.Album.Artist.Slug,
				Image: utils.ParseImageToJSON(result.Album.Artist.Image),
			},
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (svc *favoriteService) DeleteFavoriteSong(ctx context.Context, userId int, songId int) (err error) {
	exists, err := svc.favRepo.FindExistsFavoriteSongBySongId(ctx, userId, songId)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "DeleteFavoriteSong", err)
		return
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Song", Id: songId}
		utils.LogWarn(svc.log, ctx, "favorite_service", "DeleteFavoriteSong", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	if err = svc.favRepo.DeleteFavoriteSong(ctx, userId, songId); err != nil {
		utils.LogError(svc.log, ctx, "favorite_repo", "DeleteFavoriteSong", err)
		return
	}

	return
}
