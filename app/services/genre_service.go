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

func (svc *genreService) GetAll(ctx context.Context, pageSize int, offset int) (genres []dto.Genre, total int, err error) {
	// Get genres total
	total, err = svc.repo.FindCount(ctx)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetCount", err)
		return
	}

	// Get List genres
	results, err := svc.repo.FindAll(ctx, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetAll", err)
		return nil, 0, err
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

	return genres, total, nil
}

func (svc *genreService) GetGenreById(ctx context.Context, id int) (genre dto.Genre, err error) {
	result, err := svc.repo.FindGenreById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetGenreById", err)
		return
	}
	if result == nil {
		nfErr := errs.NewNotFoundError("Genre", "id", id)
		utils.LogWarn(svc.log, ctx, "genre_service", "GetGenreById", nfErr)
		return genre, nfErr
	}

	genre.Id = result.Id
	genre.Name = result.Name
	genre.Image = utils.ParseImageToJSON(result.Image)

	return genre, nil
}

func (svc *genreService) CreateGenre(ctx context.Context, req dto.CreateGenreRequest) (err error) {
	// validation failed with status bad request
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	// image image with status bad request
	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return errs.NewBadRequestError("invalid image", map[string]string{"image": "Invali image object"})
	}

	input := models.CreateGenreInput{
		Name:  req.Name,
		Image: imgByte,
	}

	return svc.repo.Store(ctx, input)
}

func (svc *genreService) UpdateGenre(ctx context.Context, req dto.CreateGenreRequest, id int) (err error) {
	if errorsMap, err := utils.RequestValidate(&req); err != nil {
		return errs.NewBadRequestError("validation failed", errorsMap)
	}

	imgByte, err := utils.ParseImageToByte(req.Image)
	if err != nil {
		return errs.NewBadRequestError("Invalid image", map[string]string{"image": "Invalid image object."})
	}

	exists, err := svc.repo.FindExistsGenreById(ctx, id)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "UpdateGenre", err)
		return err
	}

	if !exists {
		nfErr := errs.NewNotFoundError("Genre", "Id", id)
		utils.LogWarn(svc.log, ctx, "genre_repo", "UpdateGenre", nfErr)
		return nfErr
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
		nfErr := errs.NewNotFoundError("Genre", "id", id)
		utils.LogWarn(svc.log, ctx, "genre_service", "DeleteGenre", nfErr)
		return nfErr
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
		nfErr := errs.NewNotFoundError("Artist", "id", artistId)
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateArtistGenre", nfErr)
		return nfErr
	}

	// Check existing genre
	exists, err = svc.repo.FindExistsGenreById(ctx, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateArtistGenre", err)
		return err
	}
	if !exists {
		nfErr := errs.NewNotFoundError("Genre", "id", genreId)
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateArtistGenre", nfErr)
		return nfErr
	}

	// Check if genre already exists with the same artist id
	exists, err = svc.repo.FindExistsArtistGenreByGenreId(ctx, artistId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateArtistGenre", err)
		return err
	}
	if exists {
		conflictErr := errs.NewConflictError("Genre", "genre_id", genreId)
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateArtistGenre", conflictErr)
		return conflictErr
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
		nfErr := errs.NewNotFoundError("ArtistGenre", "id", artistId)
		utils.LogWarn(svc.log, ctx, "genre_service", "DeleteArtistGenre", nfErr)
		return nfErr
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
		nfErr := errs.NewNotFoundError("Song", "id", songId)
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateSongGenre", nfErr)
		return nfErr
	}

	// Check existing genre
	exists, err = svc.repo.FindExistsGenreById(ctx, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateSongGenre", err)
		return err
	}
	if !exists {
		nfErr := errs.NewNotFoundError("Genre", "id", genreId)
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateSongGenre", nfErr)
		return nfErr
	}

	// Check if genre already exists with the same song id
	exists, err = svc.repo.FindExistsSongGenreByGenreId(ctx, songId, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "CreateSongGenre", err)
		return err
	}
	if exists {
		conflictErr := errs.NewConflictError("SongGenre", "genre_id", genreId)
		utils.LogWarn(svc.log, ctx, "genre_service", "CreateSongGenre", conflictErr)
		return conflictErr
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
		nfErr := errs.NewNotFoundError("SongGenre", "id", songId)
		utils.LogWarn(svc.log, ctx, "genre_service", "DeleteSongGenre", nfErr)
		return nfErr
	}

	if err := svc.repo.DeleteSongGenre(ctx, songId, genreId); err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "DeleteSongGenre", err)
		return err
	}

	return
}

func (svc *genreService) GetAllArtists(ctx context.Context, genreId int, pageSize int, offset int) (artists []dto.Artist, total int, err error) {
	total, err = svc.repo.FindCountArtists(ctx, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetCountArtists", err)
		return nil, 0, err
	}

	results, err := svc.repo.FindAllArtists(ctx, genreId, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetAllArtists", err)
		return nil, 0, err
	}

	artists = make([]dto.Artist, 0, len(results))
	for _, result := range results {
		artists = append(artists, dto.Artist{
			Id:    result.Id,
			Name:  result.Name,
			Slug:  result.Slug,
			Image: utils.ParseImageToJSON(result.Image),
		})
	}

	return artists, total, nil
}

func (svc *genreService) GetAllSongs(ctx context.Context, genreId int, pageSize int, offset int) (songs []dto.Song, total int, err error) {
	total, err = svc.repo.FindCountSongs(ctx, genreId)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetCountSongs", err)
		return nil, 0, err
	}

	results, err := svc.repo.FindAllSongs(ctx, genreId, pageSize, offset)
	if err != nil {
		utils.LogError(svc.log, ctx, "genre_service", "GetAllSongs", err)
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
