package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type genreRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewGenreRepository(db *database.DB, log *logrus.Logger) contracts.GenreRepository {
	return &genreRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *genreRepository) FindAll(ctx context.Context, pageSize, offset int) (genres []models.Genre, err error) {
	query := `SELECT id, name, image FROM genres ORDER BY id DESC LIMIT $1 OFFSET $2`
	args := []any{pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindAll", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		genre := models.Genre{}
		if err := rows.Scan(&genre.Id, &genre.Name, &genre.Image); err != nil {
			utils.LogError(repo.log, ctx, "genre_repo", "FindAll", err)
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (repo *genreRepository) FindCount(ctx context.Context) (total int, err error) {
	query := `SELECT COUNT(*) FROM genres`

	if err = repo.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindCount", err)
		return
	}

	return
}

func (repo *genreRepository) FindExistsGenreById(ctx context.Context, id int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM genres WHERE id = $1)`

	if err = repo.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindExistsGenreById", err)
		return
	}

	return
}

func (repo *genreRepository) FindGenreById(ctx context.Context, id int) (genre *models.Genre, err error) {
	query := `SELECT id, name, image FROM genres WHERE id = $1`
	genre = &models.Genre{}

	if err := repo.db.QueryRowContext(ctx, query, id).Scan(&genre.Id, &genre.Name, &genre.Image); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarn(repo.log, ctx, "genre_repo", "FindGenreById", utils.NotFoundError{Resource: "Genre", Id: id})
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "genre_repo", "FindGenreById", err)
		return nil, err
	}

	return
}

func (repo *genreRepository) Store(ctx context.Context, input models.CreateGenreInput) (err error) {
	query := `INSERT INTO genres(name,image) VALUES($1, $2)`
	args := []any{input.Name, input.Image}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "Store", err)
		return
	}

	return
}

func (repo *genreRepository) Update(ctx context.Context, input models.CreateGenreInput, id int) (err error) {
	query := `UPDATE genres SET name = $1, image = $2 WHERE id = $3`
	args := []any{input.Name, input.Image, id}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "Update", err)
		return
	}

	return
}

func (repo *genreRepository) Delete(ctx context.Context, id int) (err error) {
	query := `DELETE FROM genres WHERE id = $1`

	if _, err = repo.db.ExecContext(ctx, query, id); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "Delete", err)
		return
	}

	return
}

func (repo *genreRepository) FindExistsArtistGenreByGenreId(ctx context.Context, artistId int, genreId int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM artist_genres WHERE artist_id = $1 AND genre_id = $2)`
	args := []any{artistId, genreId}

	if err = repo.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "FindExistsArtistGenreByGenreId", err)
		return
	}

	return
}

func (repo *genreRepository) StoreArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	query := `INSERT INTO artist_genres(artist_id, genre_id) VALUES($1, $2)`
	args := []any{artistId, genreId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "StoreArtistGenre", err)
		return
	}

	return
}

func (repo *genreRepository) FindArtistGenres(ctx context.Context, artistId, pageSize, offset int) (genres []models.Genre, err error) {
	query := `SELECT g.id AS genre_id, g.name AS genre_name, g.image AS genre_image FROM artist_genres ag INNER JOIN genres g ON g.id = ag.genre_id WHERE ag.artist_id = $1 ORDER BY ag.created_at DESC LIMIT $2 OFFSET $3`
	args := []any{artistId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "FindArtistGenres", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.Id, &genre.Name, &genre.Image); err != nil {
			utils.LogError(repo.log, ctx, "artist_repo", "FindArtistGenres", err)
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (repo *genreRepository) DeleteArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	query := `DELETE FROM artist_genres WHERE artist_id = $1 AND genre_id = $2`
	args := []any{artistId, genreId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "DeleteArtistGenre", err)
		return err
	}

	return
}
