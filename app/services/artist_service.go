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

type artistService struct {
	repo contracts.ArtistRepository
	log  *logrus.Logger
}

func NewArtistService(repo contracts.ArtistRepository, log *logrus.Logger) contracts.ArtistService {
	return &artistService{
		repo: repo,
		log:  log,
	}
}

func (svc *artistService) GetAll(ctx context.Context, pageSize, offset int) (artists []dto.Artist, total int, err error) {
	total, err = svc.repo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "GetAll", err)
		return nil, 0, err
	}

	results, err := svc.repo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "GetAll", err)
		return nil, 0, err
	}

	artists = make([]dto.Artist, 0, len(results))
	for _, result := range results {
		artist := dto.Artist{
			Id:    result.Id,
			Name:  result.Name,
			Slug:  result.Slug,
			Image: utils.ParseImageToJSON(result.Image),
		}
		artists = append(artists, artist)
	}

	return artists, total, nil
}

func (svc *artistService) CreateArtist(ctx context.Context, req dto.CreateArtistRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	slug := utils.MakeSlug(req.Name)
	exists, err := svc.repo.FindExistsArtistBySlug(ctx, slug)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "CreateArtist", err)
		return err
	}
	if exists {
		conflictErr := errs.NewConflictError("Artist", "name", req.Name)
		utils.LogWarn(svc.log, ctx, "artist_service", "CreateArtist", conflictErr)
		return conflictErr
	}

	input := models.CreateArtistInput{
		Name:  req.Name,
		Slug:  slug,
		Image: utils.ParseImageToByte(req.Image),
	}

	if err := svc.repo.Store(ctx, input); err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "CreateArtist", err)
		return err
	}

	return
}

func (svc *artistService) GetArtistById(ctx context.Context, artistId int) (artist dto.Artist, err error) {
	result, err := svc.repo.FindArtistById(ctx, artistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "GetArtistById", err)
		return
	}
	if result == nil {
		notFoundErr := errs.NewNotFoundError("Artist", "id", artistId)
		utils.LogWarn(svc.log, ctx, "artist_service", "GetArtistById", notFoundErr)
		return artist, notFoundErr
	}

	artist = dto.Artist{
		Id:    result.Id,
		Name:  result.Name,
		Slug:  result.Slug,
		Image: utils.ParseImageToJSON(result.Image),
	}

	return artist, nil
}

func (svc *artistService) UpdateArtist(ctx context.Context, req dto.CreateArtistRequest, id int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	artist, err := svc.repo.FindArtistById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "UpdateArtist", err)
		return err
	}
	if artist == nil {
		notFoundErr := errs.NewNotFoundError("Artist", "id", id)
		utils.LogWarn(svc.log, ctx, "artist_service", "UpdateArtist", notFoundErr)
		return notFoundErr
	}

	slug := utils.MakeSlug(req.Name)
	if artist.Slug != slug {
		exists, err := svc.repo.FindExistsArtistBySlug(ctx, slug)
		if err != nil {
			utils.LogError(svc.log, ctx, "artist_service", "UpdateArtist", err)
			return err
		}
		if exists {
			conflictErr := errs.NewConflictError("Artist", "name", req.Name)
			utils.LogWarn(svc.log, ctx, "artist_service", "UpdateArtist", conflictErr)
			return conflictErr
		}
	}

	input := models.CreateArtistInput{
		Name:  req.Name,
		Slug:  slug,
		Image: utils.ParseImageToByte(req.Image),
	}

	if err := svc.repo.Update(ctx, input, id); err != nil {
		utils.LogWarn(svc.log, ctx, "artist_service", "UpdateArtist", err)
		return err
	}

	return
}

func (svc *artistService) DeleteArtist(ctx context.Context, id int) (err error) {
	exists, err := svc.repo.FindExistsArtistById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "DeleteArtist", err)
		return err
	}
	if !exists {
		notFoundErr := errs.NewNotFoundError("Artist", "id", id)
		utils.LogWarn(svc.log, ctx, "artist_service", "DeleteArtist", notFoundErr)
		return notFoundErr
	}

	if err := svc.repo.Delete(ctx, id); err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "DeleteArtist", err)
		return err
	}

	return
}
