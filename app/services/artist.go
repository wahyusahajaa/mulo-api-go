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

type artistService struct {
	repo      contracts.ArtistRepository
	genreRepo contracts.GenreRepository
	log       *logrus.Logger
}

func NewArtistService(repo contracts.ArtistRepository, genreRepo contracts.GenreRepository, log *logrus.Logger) contracts.ArtistService {
	return &artistService{
		repo:      repo,
		genreRepo: genreRepo,
		log:       log,
	}
}

func (svc *artistService) GetAll(ctx context.Context, pageSize, offset int) (artists []dto.Artist, err error) {
	results, err := svc.repo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "GetAll", err)
		return nil, err
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

	return artists, nil
}

func (svc *artistService) GetArtistByIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error) {
	artists, err = svc.repo.FindByArtistIds(ctx, inClause, artistIds)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}

func (svc *artistService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "GetCount", err)
		return
	}

	return
}

func (svc *artistService) CreateArtist(ctx context.Context, req dto.CreateArtistRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}
	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	slug := utils.MakeSlug(req.Name)
	exists, err := svc.repo.FindExistsArtistBySlug(ctx, slug)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "CreateArtist", err)
		return err
	}
	if exists {
		conflictErr := utils.ConflictError{Resource: "Artist", Field: "name", Value: req.Name}
		utils.LogWarn(svc.log, ctx, "artist_service", "CreateArtist", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	input := models.CreateArtistInput{
		Name:  req.Name,
		Slug:  slug,
		Image: imgByte,
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
		notFoundErr := utils.NotFoundError{Resource: "Artist", Id: artistId}
		utils.LogWarn(svc.log, ctx, "artist_service", "GetArtistById", notFoundErr)
		return artist, fmt.Errorf("%w", notFoundErr)
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
		return fmt.Errorf("%w", utils.BadReqError{Errors: errorsMap})
	}

	artist, err := svc.repo.FindArtistById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "UpdateArtist", err)
		return err
	}
	if artist == nil {
		notFoundErr := utils.NotFoundError{Resource: "Artist", Id: id}
		utils.LogWarn(svc.log, ctx, "artist_service", "UpdateArtist", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	slug := utils.MakeSlug(req.Name)
	if artist.Slug != slug {
		exists, err := svc.repo.FindExistsArtistBySlug(ctx, slug)
		if err != nil {
			utils.LogError(svc.log, ctx, "artist_service", "UpdateArtist", err)
			return err
		}
		if exists {
			conflictErr := utils.ConflictError{Resource: "Artist", Field: "Name", Value: req.Name}
			utils.LogWarn(svc.log, ctx, "artist_service", "UpdateArtist", conflictErr)
			return fmt.Errorf("%w", conflictErr)
		}
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("%w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	input := models.CreateArtistInput{
		Name:  req.Name,
		Slug:  slug,
		Image: imgByte,
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
		notFoundErr := utils.NotFoundError{Resource: "Artist", Id: id}
		utils.LogWarn(svc.log, ctx, "artist_service", "DeleteArtist", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	if err := svc.repo.Delete(ctx, id); err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "DeleteArtist", err)
		return err
	}

	return
}

func (svc *artistService) CreateArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	// Check existing artist
	exists, err := svc.repo.FindExistsArtistById(ctx, artistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "CreateArtistGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Artist", Id: artistId}
		utils.LogWarn(svc.log, ctx, "artist_service", "CreateArtistGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check existing genre
	exists, err = svc.genreRepo.FindExistsGenreById(ctx, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "CreateArtistGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Genre", Id: genreId}
		utils.LogWarn(svc.log, ctx, "artist_service", "CreateArtistGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check if genre already exists with the same artist id
	exists, err = svc.repo.FindExistsArtistGenreByGenreId(ctx, artistId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "CreateArtistGenre", err)
		return err
	}
	if exists {
		conflictErr := utils.ConflictError{Resource: "Genre", Field: "genre_id", Value: genreId}
		utils.LogWarn(svc.log, ctx, "artist_service", "CreateArtistGenre", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	// Store new artist genre
	if err = svc.repo.StoreArtistGenre(ctx, artistId, genreId); err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "CreateArtistGenre", err)
		return err
	}

	return
}

func (svc *artistService) GetArtistGenres(ctx context.Context, artistId int, pageSize int, offset int) (genres []dto.Genre, err error) {
	results, err := svc.repo.FindArtistGenres(ctx, artistId, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "GetArtistGenres", err)
		return genres, err
	}

	genres = make([]dto.Genre, 0, len(results))
	for _, result := range results {
		genre := dto.Genre{
			Id:    result.Id,
			Name:  result.Name,
			Image: utils.ParseImageToJSON(result.Image),
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (svc *artistService) DeleteArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	exists, err := svc.repo.FindExistsArtistGenreByGenreId(ctx, artistId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "DeleteArtistGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "ArtistGenre", Id: artistId}
		utils.LogWarn(svc.log, ctx, "artist_service", "DeleteArtistGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	if err := svc.repo.DeleteArtistGenre(ctx, artistId, genreId); err != nil {
		utils.LogError(svc.log, ctx, "artist_service", "DeleteArtistGenre", err)
		return err
	}

	return
}
