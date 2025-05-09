package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type artistRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewArtistRepository(db *database.DB, log *logrus.Logger) contracts.ArtistRepository {
	return &artistRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *artistRepository) FindAll(ctx context.Context, pageSize, offset int) (artists []models.Artist, err error) {
	query := `SELECT id, name, slug, image FROM artists ORDER BY id DESC LIMIT $1 OFFSET $2`
	args := []any{pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "FindAll", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		artist := models.Artist{}
		if err := rows.Scan(
			&artist.Id,
			&artist.Name,
			&artist.Slug,
			&artist.Image,
		); err != nil {
			utils.LogError(repo.log, ctx, "artist_repo", "FindAll", err)
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func (repo *artistRepository) FindByArtistIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error) {
	query := fmt.Sprintf(`SELECT id, name, slug, image FROM artists WHERE id IN %s`, inClause)

	rows, err := repo.db.QueryContext(ctx, query, artistIds...)
	if err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "FindByArtistIds", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		artist := models.Artist{}
		if err := rows.Scan(
			&artist.Id,
			&artist.Name,
			&artist.Slug,
			&artist.Image,
		); err != nil {
			utils.LogError(repo.log, ctx, "artist_repo", "FindByArtistIds", err)
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func (repo *artistRepository) FindExistsArtistBySlug(ctx context.Context, slug string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM artists WHERE slug = $1)`

	if err = repo.db.QueryRowContext(ctx, query, slug).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "FindDuplicateArtistBySlug", err)
		return
	}

	return
}

func (repo *artistRepository) FindExistsArtistById(ctx context.Context, id int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM artists WHERE id = $1)`

	if err = repo.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "FindExistsArtistById", err)
		return
	}

	return
}

func (repo *artistRepository) FindCount(ctx context.Context) (total int, err error) {
	query := `SELECT COUNT(*) FROM artists`

	if err := repo.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "Count", err)
		return 0, err
	}

	return
}

func (repo *artistRepository) Store(ctx context.Context, input models.CreateArtistInput) (err error) {
	query := `INSERT INTO artists(name, slug, image) VALUES($1, $2, $3)`
	args := []any{input.Name, input.Slug, input.Image}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "Store", err)
		return
	}

	return
}

func (repo *artistRepository) FindArtistById(ctx context.Context, artistId int) (artist *models.Artist, err error) {
	query := `SELECT id, name, slug, image FROM artists WHERE id = $1`
	artist = &models.Artist{}

	if err = repo.db.QueryRowContext(ctx, query, artistId).Scan(
		&artist.Id,
		&artist.Name,
		&artist.Slug,
		&artist.Image,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFoundErr := utils.NotFoundError{Resource: "Artist", Id: artistId}
			utils.LogWarn(repo.log, ctx, "artist_repo", "FindArtistById", notFoundErr)
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "artist_repo", "FindArtistById", err)
		return nil, err
	}

	return artist, nil
}

func (repo *artistRepository) Update(ctx context.Context, input models.CreateArtistInput, id int) (err error) {
	query := `UPDATE artists SET name = $1, slug = $2, image = $3 WHERE id = $4`
	args := []any{input.Name, input.Slug, input.Image, id}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "Update", err)
		return err
	}

	return
}

func (repo *artistRepository) Delete(ctx context.Context, id int) (err error) {
	query := `DELETE FROM artists WHERE id = $1`

	if _, err = repo.db.ExecContext(ctx, query, id); err != nil {
		utils.LogError(repo.log, ctx, "artist_repo", "Delete", err)
		return
	}

	return
}
