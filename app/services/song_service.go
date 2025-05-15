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

type songService struct {
	songRepo  contracts.SongRepository
	albumRepo contracts.AlbumRepository
	log       *logrus.Logger
}

func NewSongService(songRepo contracts.SongRepository, albumRepo contracts.AlbumRepository, log *logrus.Logger) contracts.SongService {
	return &songService{
		songRepo:  songRepo,
		albumRepo: albumRepo,
		log:       log,
	}
}

func (svc *songService) GetAll(ctx context.Context, pageSize int, offset int) (songs []dto.Song, total int, err error) {
	total, err = svc.songRepo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "GetAll", err)
		return nil, 0, err
	}

	results, err := svc.songRepo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "GetAll", err)
		return nil, 0, err
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

	return songs, total, nil
}

func (svc *songService) GetSongById(ctx context.Context, id int) (song dto.Song, err error) {
	result, err := svc.songRepo.FindSongById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "GetSongById", err)
		return
	}

	if result == nil {
		notFoundErr := errs.NewNotFoundError("Song", "id", id)
		utils.LogWarn(svc.log, ctx, "song_service", "GetSongById", notFoundErr)
		return song, notFoundErr
	}

	song.Id = result.Id
	song.Title = result.Title
	song.Audio = result.Audio
	song.Duration = result.Duration
	song.Image = utils.ParseImageToJSON(result.Image)
	song.Album = dto.AlbumWithArtist{
		Album: dto.Album{
			Id:    result.Album.Id,
			Name:  result.Album.Name,
			Slug:  result.Album.Slug,
			Image: utils.ParseImageToJSON(result.Album.Artist.Image),
		},
		Artist: dto.Artist{
			Id:    result.Album.Artist.Id,
			Name:  result.Album.Artist.Name,
			Slug:  result.Album.Artist.Slug,
			Image: utils.ParseImageToJSON(result.Album.Artist.Image),
		},
	}

	return
}

func (svc *songService) CreateSong(ctx context.Context, req dto.CreateSongRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	exists, err := svc.albumRepo.FindExistsAlbumById(ctx, req.AlbumId)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "CreateSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Album", "id", req.AlbumId)
		utils.LogWarn(svc.log, ctx, "song_service", "CreateSong", notFoundErr)
		return notFoundErr
	}

	input := models.CreateSongInput{
		AlbumId:  req.AlbumId,
		Title:    req.Title,
		Audio:    req.Audio,
		Duration: req.Duration,
		Image:    utils.ParseImageToByte(req.Image),
	}

	if err := svc.songRepo.Store(ctx, input); err != nil {
		utils.LogError(svc.log, ctx, "song_service", "CreateSong", err)
		return err
	}

	return
}

func (svc *songService) UpdateSong(ctx context.Context, req dto.CreateSongRequest, id int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	exists, err := svc.songRepo.FindExistsSongById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "UpdateSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Song", "id", id)
		utils.LogError(svc.log, ctx, "song_service", "UpdateSong", notFoundErr)
		return notFoundErr
	}

	exists, err = svc.albumRepo.FindExistsAlbumById(ctx, req.AlbumId)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "UpdateSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Album", "id", req.AlbumId)
		utils.LogWarn(svc.log, ctx, "song_service", "UpdateSong", notFoundErr)
		return notFoundErr
	}

	input := models.CreateSongInput{
		AlbumId:  req.AlbumId,
		Title:    req.Title,
		Audio:    req.Audio,
		Duration: req.Duration,
		Image:    utils.ParseImageToByte(req.Image),
	}

	if err := svc.songRepo.Update(ctx, input, id); err != nil {
		utils.LogError(svc.log, ctx, "song_service", "UpdateSong", err)
		return err
	}

	return
}

func (svc *songService) DeleteSong(ctx context.Context, id int) (err error) {
	exists, err := svc.songRepo.FindExistsSongById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "DeleteSong", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Song", "id", id)
		utils.LogWarn(svc.log, ctx, "song_service", "DeleteSong", notFoundErr)
		return notFoundErr
	}

	if err := svc.songRepo.Delete(ctx, id); err != nil {
		utils.LogError(svc.log, ctx, "song_service", "DeleteSong", err)
		return err
	}

	return
}
