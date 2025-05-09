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

type genreService struct {
	repo       contracts.GenreRepository
	artistRepo contracts.ArtistRepository
	songRepo   contracts.SongRepository
	log        *logrus.Logger
}

func NewGenreService(repo contracts.GenreRepository, artistRepo contracts.ArtistRepository, songRepo contracts.SongRepository, log *logrus.Logger) contracts.GenreService {
	return &genreService{
		repo:       repo,
		artistRepo: artistRepo,
		songRepo:   songRepo,
		log:        log,
	}
}

func (svc *genreService) GetAll(ctx context.Context, pageSize int, offset int) (genres []dto.Genre, err error) {
	results, err := svc.repo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetAll", err)
		return nil, err
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

func (svc *genreService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetCount", err)
		return
	}

	return
}

func (svc *genreService) GetGenreById(ctx context.Context, id int) (genre dto.Genre, err error) {
	result, err := svc.repo.FindGenreById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetGenreById", err)
		return
	}
	if result == nil {
		notFoundErr := utils.NotFoundError{Resource: "Genre", Id: id}
		utils.LogWarn(svc.log, ctx, "genre_service", "GetGenreById", notFoundErr)
		return genre, fmt.Errorf("not_found: %w", notFoundErr)
	}

	genre.Id = result.Id
	genre.Name = result.Name
	genre.Image = utils.ParseImageToJSON(result.Image)

	return genre, nil
}

func (svc *genreService) CreateGenre(ctx context.Context, req dto.CreateGenreRequest) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("validation: %w", utils.BadReqError{Errors: errorsMap})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("parse_image: %w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	input := models.CreateGenreInput{
		Name:  req.Name,
		Image: imgByte,
	}

	return svc.repo.Store(ctx, input)
}

func (svc *genreService) UpdateGenre(ctx context.Context, req dto.CreateGenreRequest, id int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return fmt.Errorf("validation: %w", utils.BadReqError{Errors: errorsMap})
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return fmt.Errorf("parse_image: %w", utils.BadReqError{Errors: map[string]string{"image": "Invalid image object"}})
	}

	exists, err := svc.repo.FindExistsGenreById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "UpdateGenre", err)
		return err
	}

	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Genre", Id: id}
		utils.LogWarn(svc.log, ctx, "genre_repo", "UpdateGenre", notFoundErr)
		return fmt.Errorf("not_found: %w", notFoundErr)
	}

	input := models.CreateGenreInput{
		Name:  req.Name,
		Image: imgByte,
	}

	if err := svc.repo.Update(ctx, input, id); err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "UpdateGenre", err)
		return err
	}

	return
}

func (svc *genreService) DeleteGenre(ctx context.Context, id int) (err error) {
	exists, err := svc.repo.FindExistsGenreById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "DeleteGenre", err)
		return
	}

	if !exists {
		errNotFound := utils.NotFoundError{Resource: "Genre", Id: id}
		utils.LogWarn(svc.log, ctx, "genre_service", "DeleteGenre", errNotFound)
		return fmt.Errorf("not_found: %w", errNotFound)
	}

	if err = svc.repo.Delete(ctx, id); err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "DeleteGenre", err)
		return
	}

	return
}

func (svc *genreService) CreateArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	// Check existing artist
	exists, err := svc.artistRepo.FindExistsArtistById(ctx, artistId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateArtistGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Artist", Id: artistId}
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateArtistGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check existing genre
	exists, err = svc.repo.FindExistsGenreById(ctx, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateArtistGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Genre", Id: genreId}
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateArtistGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check if genre already exists with the same artist id
	exists, err = svc.repo.FindExistsArtistGenreByGenreId(ctx, artistId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateArtistGenre", err)
		return err
	}
	if exists {
		conflictErr := utils.ConflictError{Resource: "Genre", Field: "genre_id", Value: genreId}
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateArtistGenre", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	// Store new artist genre
	if err = svc.repo.StoreArtistGenre(ctx, artistId, genreId); err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateArtistGenre", err)
		return err
	}

	return
}

func (svc *genreService) GetArtistGenres(ctx context.Context, artistId int, pageSize int, offset int) (genres []dto.Genre, err error) {
	results, err := svc.repo.FindArtistGenres(ctx, artistId, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetArtistGenres", err)
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

func (svc *genreService) DeleteArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	exists, err := svc.repo.FindExistsArtistGenreByGenreId(ctx, artistId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "DeleteArtistGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "ArtistGenre", Id: artistId}
		utils.LogWarn(svc.log, ctx, "genre_service", "DeleteArtistGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	if err := svc.repo.DeleteArtistGenre(ctx, artistId, genreId); err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "DeleteArtistGenre", err)
		return err
	}

	return
}

func (svc *genreService) CreateSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	// Check existing song
	exists, err := svc.songRepo.FindExistsSongById(ctx, songId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateSongGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Song", Id: songId}
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateSongGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check existing genre
	exists, err = svc.repo.FindExistsGenreById(ctx, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateSongGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "Genre", Id: genreId}
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateSongGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	// Check if genre already exists with the same song id
	exists, err = svc.repo.FindExistsSongGenreByGenreId(ctx, songId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateSongGenre", err)
		return err
	}
	if exists {
		conflictErr := utils.ConflictError{Resource: "Genre", Field: "genre_id", Value: genreId}
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateSongGenre", conflictErr)
		return fmt.Errorf("%w", conflictErr)
	}

	// Store new song genre
	if err = svc.repo.StoreSongGenre(ctx, songId, genreId); err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateSongGenre", err)
		return err
	}

	return
}

func (svc *genreService) GetSongGenres(ctx context.Context, songId int, pageSize int, offset int) (genres []dto.Genre, err error) {
	results, err := svc.repo.FindSongGenres(ctx, songId, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetSongGenres", err)
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

func (svc *genreService) DeleteSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	exists, err := svc.repo.FindExistsSongGenreByGenreId(ctx, songId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "DeleteSongGenre", err)
		return err
	}
	if !exists {
		notFoundErr := utils.NotFoundError{Resource: "SongGenre", Id: songId}
		utils.LogWarn(svc.log, ctx, "genre_service", "DeleteSongGenre", notFoundErr)
		return fmt.Errorf("%w", notFoundErr)
	}

	if err := svc.repo.DeleteSongGenre(ctx, songId, genreId); err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "DeleteSongGenre", err)
		return err
	}

	return
}
