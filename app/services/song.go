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

func (svc *songService) GetAll(ctx context.Context, pageSize int, offset int) (songs []dto.Song, err error) {
	results, err := svc.songRepo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "GetAll", err)
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
			Album: dto.Album{
				Id:    v.Album.Id,
				Name:  v.Album.Name,
				Slug:  v.Album.Slug,
				Image: utils.ParseImageToJSON(v.Album.Image),
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

func (svc *songService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.songRepo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "GetCount", err)
		return
	}
	return
}

func (svc *songService) GetSongById(ctx context.Context, id int) (song dto.Song, err error) {
	result, err := svc.songRepo.FindSongById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "GetSongById", err)
		return
	}

	if result == nil {
		notFoundErr := utils.NotFoundError{Resource: "Song", Id: id}
		utils.LogWarn(svc.log, ctx, "song_service", "GetSongById", notFoundErr)
		return song, fmt.Errorf("not_found: %w", notFoundErr)
	}

	song.Id = result.Id
	song.Title = result.Title
	song.Audio = result.Audio
	song.Duration = result.Duration
	song.Image = utils.ParseImageToJSON(result.Image)
	song.Album = dto.Album{
		Id:    result.Album.Id,
		Name:  result.Album.Name,
		Slug:  result.Album.Slug,
		Image: utils.ParseImageToJSON(result.Album.Artist.Image),
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
		return fmt.Errorf("validation: %w", utils.BadReqError{Errors: errorsMap})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("transform: %w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	exists, err := svc.albumRepo.FindExistsAlbumById(ctx, req.AlbumId)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "CreateSong", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Album", Id: req.AlbumId}
		utils.LogWarn(svc.log, ctx, "song_service", "CreateSong", notFoundErr)
		return fmt.Errorf("not_found: %w", notFoundErr)
	}

	input := models.CreateSongInput{
		AlbumId:  req.AlbumId,
		Title:    req.Title,
		Audio:    req.Audio,
		Duration: req.Duration,
		Image:    imgByte,
	}

	if err := svc.songRepo.Store(ctx, input); err != nil {
		utils.LogError(svc.log, ctx, "song_service", "CreateSong", err)
		return err
	}

	return
}

func (svc *songService) UpdateSong(ctx context.Context, req dto.CreateSongRequest, id int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("validation: %w", utils.BadReqError{Errors: errorsMap})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("transform: %w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	exists, err := svc.songRepo.FindExistsSongById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "UpdateSong", err)
		return err
	}
	if !exists {
		return fmt.Errorf("not_found: %w", utils.NotFoundError{Resource: "Song", Id: id})
	}

	exists, err = svc.albumRepo.FindExistsAlbumById(ctx, req.AlbumId)
	if err != nil {
		utils.LogError(svc.log, ctx, "song_service", "UpdateSong", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Album", Id: req.AlbumId}
		utils.LogWarn(svc.log, ctx, "song_service", "UpdateSong", notFoundErr)
		return fmt.Errorf("not_found: %w", notFoundErr)
	}

	input := models.CreateSongInput{
		AlbumId:  req.AlbumId,
		Title:    req.Title,
		Audio:    req.Audio,
		Duration: req.Duration,
		Image:    imgByte,
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
		notFoundErr := utils.NotFoundError{Resource: "Song", Id: id}
		utils.LogWarn(svc.log, ctx, "song_service", "DeleteSong", notFoundErr)
		return fmt.Errorf("not_found: %w", notFoundErr)
	}

	if err := svc.songRepo.Delete(ctx, id); err != nil {
		utils.LogError(svc.log, ctx, "song_service", "DeleteSong", err)
		return err
	}

	return
}
