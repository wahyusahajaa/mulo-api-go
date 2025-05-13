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

type albumService struct {
	repo       contracts.AlbumRepository
	artistRepo contracts.ArtistRepository
	log        *logrus.Logger
}

func NewAlbumService(repo contracts.AlbumRepository, artistRepo contracts.ArtistRepository, log *logrus.Logger) contracts.AlbumService {
	return &albumService{
		repo:       repo,
		artistRepo: artistRepo,
		log:        log,
	}
}

func (svc *albumService) GetAll(ctx context.Context, pageSize int, offset int) (albums []dto.AlbumWithArtist, err error) {
	albumResults, err := svc.repo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "GetAll", err)
		return nil, err
	}

	// Make unique artist ids
	artistIdMap := make(map[int]struct{})
	for _, albumResult := range albumResults {
		artistIdMap[albumResult.ArtistId] = struct{}{}
	}
	artistIds := make([]any, 0, len(artistIdMap))
	for id := range artistIdMap {
		artistIds = append(artistIds, id)
	}
	inClause, args := utils.BuildInClause(1, artistIds)

	// Get artists by artist ids
	artists, err := svc.artistRepo.FindByArtistIds(ctx, inClause, args)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "GetAll", err)
		return nil, err
	}

	//  Build artist lookup map
	artistMap := make(map[int]models.Artist)
	for _, artist := range artists {
		artistMap[artist.Id] = artist
	}

	albums = make([]dto.AlbumWithArtist, 0, len(albumResults))
	for _, albumResult := range albumResults {
		album := dto.AlbumWithArtist{
			Album: dto.Album{
				Id:    albumResult.Id,
				Name:  albumResult.Name,
				Slug:  albumResult.Slug,
				Image: utils.ParseImageToJSON(albumResult.Image),
			},
		}

		if artist, ok := artistMap[albumResult.ArtistId]; ok {
			album.Artist.Id = artist.Id
			album.Artist.Name = artist.Name
			album.Artist.Slug = artist.Slug
			album.Artist.Image = utils.ParseImageToJSON(artist.Image)
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func (svc *albumService) GetAlbumById(ctx context.Context, id int) (album dto.AlbumWithArtist, err error) {
	result, err := svc.repo.FindAlbumById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "GetAlbumById", err)
		return album, err
	}
	if result == nil {
		notFoundErr := utils.NotFoundError{Resource: "Album", Id: id}
		utils.LogWarn(svc.log, ctx, "album_service", "GetAlbumById", notFoundErr)
		return album, fmt.Errorf("%w", notFoundErr)
	}

	album = dto.AlbumWithArtist{
		Album: dto.Album{
			Id:    result.Id,
			Name:  result.Name,
			Slug:  result.Slug,
			Image: utils.ParseImageToJSON(result.Image),
		},
		Artist: dto.Artist{
			Id:    result.Artist.Id,
			Name:  result.Artist.Name,
			Slug:  result.Artist.Slug,
			Image: utils.ParseImageToJSON(result.Artist.Image),
		},
	}

	return album, nil
}

func (svc *albumService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "GetCount", err)
		return
	}

	return
}

func (svc *albumService) CreateAlbum(ctx context.Context, req dto.CreateAlbumRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	// Check existing artist
	exists, err := svc.artistRepo.FindExistsArtistById(ctx, req.ArtistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "CreateAlbum", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Artist", Id: req.ArtistId}
		utils.LogWarn(svc.log, ctx, "album_service", "CreateAlbum", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check duplicate slug
	slug := utils.MakeSlug(req.Name)
	exist, err := svc.repo.FindExistsAlbumBySlug(ctx, slug)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "CreateAlbum", err)
		return err
	}
	if exist {
		conflictErr := utils.ConflictError{Resource: "Album", Field: "Name", Value: req.Name}
		utils.LogWarn(svc.log, ctx, "album_service", "CreateAlbum", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	input := models.CreateAlbumInput{
		ArtistId: req.ArtistId,
		Name:     req.Name,
		Slug:     slug,
		Image:    imgByte,
	}

	if err := svc.repo.Store(ctx, input); err != nil {
		utils.LogError(svc.log, ctx, "album_service", "CreateAlbum", err)
		return err
	}

	return
}

func (svc *albumService) UpdateAlbum(ctx context.Context, req dto.CreateAlbumRequest, id int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	album, err := svc.repo.FindAlbumById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "UpdateAlbum", err)
		return err
	}
	if album == nil {
		notFoundErr := utils.NotFoundError{Resource: "Album", Id: id}
		utils.LogWarn(svc.log, ctx, "album_service", "UpdateAlbum", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check existing artist
	exists, err := svc.artistRepo.FindExistsArtistById(ctx, req.ArtistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "UpdateAlbum", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Artist", Id: req.ArtistId}
		utils.LogWarn(svc.log, ctx, "album_service", "UpdateAlbum", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	slug := utils.MakeSlug(req.Name)
	if album.Slug != slug {
		exists, err := svc.repo.FindExistsAlbumBySlug(ctx, slug)
		if err != nil {
			utils.LogError(svc.log, ctx, "album_service", "UpdateAlbum", err)
			return err
		}
		if exists {
			conflictErr := utils.ConflictError{Resource: "Album", Field: "Name", Value: req.Name}
			utils.LogWarn(svc.log, ctx, "album_service", "UpdateAlbum", conflictErr)
			return fmt.Errorf("%w", conflictErr)
		}
	}

	input := models.CreateAlbumInput{
		ArtistId: req.ArtistId,
		Name:     req.Name,
		Slug:     slug,
		Image:    imgByte,
	}

	if err := svc.repo.Update(ctx, input, id); err != nil {
		utils.LogError(svc.log, ctx, "album_service", "UpdateAlbum", err)
		return err
	}

	return
}

func (svc *albumService) DeleteAlbum(ctx context.Context, id int) (err error) {
	exists, err := svc.repo.FindExistsAlbumById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "DeleteAlbum", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Album", Id: id}
		utils.LogWarn(svc.log, ctx, "album_service", "DeleteAlbum", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	if err = svc.repo.Delete(ctx, id); err != nil {
		utils.LogError(svc.log, ctx, "album_service", "DeleteAlbum", err)
		return
	}

	return
}

func (svc *albumService) GetAlbumsByArtistId(ctx context.Context, artistId int) (albums []dto.Album, err error) {
	results, err := svc.repo.FindAlbumsByArtistId(ctx, artistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "album_service", "GetAlbumsByArtistId", err)
		return nil, err
	}

	albums = make([]dto.Album, 0, len(results))
	for _, result := range results {
		album := dto.Album{
			Id:    result.Id,
			Name:  result.Name,
			Slug:  result.Slug,
			Image: utils.ParseImageToJSON(result.Image),
		}

		albums = append(albums, album)
	}

	return albums, nil
}
