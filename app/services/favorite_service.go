package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
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

func (svc *favoriteService) GetFavoriteSongsByUserID(ctx context.Context, userID, pageSize, offset int) (songs []dto.Song, total int, err error) {
	total, err = svc.favRepo.FindCountFavoriteSongsByUserID(ctx, userID)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "GetFavoriteSongsByUserID", err)
		return nil, 0, err
	}

	results, err := svc.favRepo.FindFavoriteSongsByUserID(ctx, userID, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "GetFavoriteSongsByUserID", err)
		return nil, 0, err
	}

	songs = make([]dto.Song, 0, len(results))
	for _, result := range results {
		var song dto.Song
		song.Id = result.Id
		song.Title = result.Title
		song.Duration = result.Duration
		song.Audio = result.Audio
		song.Image = utils.ParseImageToJSON(result.Image)
		song.Album = dto.AlbumWithArtist{
			Album: dto.Album{
				Id:    result.Album.Id,
				Name:  result.Album.Name,
				Slug:  result.Album.Slug,
				Image: utils.ParseImageToJSON(result.Album.Image),
			},
			Artist: dto.Artist{
				Id:    result.Album.Artist.Id,
				Name:  result.Album.Artist.Name,
				Slug:  result.Album.Artist.Slug,
				Image: utils.ParseImageToJSON(result.Album.Artist.Image),
			},
		}

		songs = append(songs, song)
	}

	return songs, total, nil
}

func (svc *favoriteService) AddFavoriteSong(ctx context.Context, userID int, songID int) (err error) {
	exists, err := svc.songRepo.FindExistsSongById(ctx, songID)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "CreateFavorite", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Song", "id", songID)
		utils.LogWarn(svc.log, ctx, "favorite_service", "CreateFavorite", notFoundErr)
		return notFoundErr
	}

	exists, err = svc.favRepo.FindExistsFavoriteSongBySongID(ctx, userID, songID)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "CreateFavorite", err)
		return err
	}

	if exists {
		conflictErr := errs.NewConflictError("Favorite", "song_id", songID)
		utils.LogWarn(svc.log, ctx, "favorite_service", "CreateFavorite", conflictErr)
		return conflictErr
	}

	if err = svc.favRepo.StoreFavoriteSong(ctx, userID, songID); err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "CreateFavorite", err)
		return err
	}

	return
}

func (svc *favoriteService) RemoveFavoriteSong(ctx context.Context, userID int, songID int) (err error) {
	exists, err := svc.favRepo.FindExistsFavoriteSongBySongID(ctx, userID, songID)
	if err != nil {
		utils.LogError(svc.log, ctx, "favorite_service", "DeleteFavoriteSong", err)
		return
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Song", "id", songID)
		utils.LogWarn(svc.log, ctx, "favorite_service", "DeleteFavoriteSong", notFoundErr)
		return notFoundErr
	}

	if err = svc.favRepo.DeleteFavoriteSong(ctx, userID, songID); err != nil {
		utils.LogError(svc.log, ctx, "favorite_repo", "DeleteFavoriteSong", err)
		return
	}

	return
}
