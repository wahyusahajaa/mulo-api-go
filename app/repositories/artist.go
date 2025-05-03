package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
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

func (repo *artistRepository) FetchAll(ctx context.Context, pageSize, offset int) (artists []models.Artist, err error) {
	query := `SELECT id, name, slug, image FROM artists ORDER BY id DESC LIMIT $1 OFFSET $2`
	args := []any{pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)

	if err != nil {
		repo.log.WithError(err).Error("failed to query artists")
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
			repo.log.WithError(err).Error("failed to scan artist")
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func (repo *artistRepository) Count(ctx context.Context) (total int, err error) {
	query := `SELECT COUNT(*) FROM artists`

	if err := repo.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		repo.log.WithError(err).Error("failed to query count artist")
		return 0, err
	}

	return
}

func (repo *artistRepository) Store(ctx context.Context, name, slug string, image []byte) (err error) {
	query := `INSERT INTO artists(name, slug, image) VALUES($1, $2, $3)`
	args := []any{name, slug, image}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.log.WithError(err).Error("failed to query insert artists")
		return
	}

	return
}

func (repo *artistRepository) FindDuplicateArtistBySlug(ctx context.Context, slug string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM artists WHERE slug = $1)`

	if err = repo.db.QueryRowContext(ctx, query, slug).Scan(&exists); err != nil {
		repo.log.WithError(err).Error("failed to query check duplicate artist")
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
			repo.log.WithField("artist_id", artistId).Warn("artist not found")
			return nil, nil
		}

		repo.log.WithError(err).Error("failed to query artist by id")
		return nil, err
	}

	return artist, nil
}

func (repo *artistRepository) Update(ctx context.Context, name string, slug string, image []byte, artistId int) (err error) {
	query := `UPDATE artists SET name = $1, slug = $2, image = $3 WHERE id = $4`
	args := []any{name, slug, image, artistId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		repo.log.WithError(err).Error("failed to query update artist")
		return err
	}

	return
}

func (repo *artistRepository) Delete(ctx context.Context, artistId int) (err error) {
	query := `DELETE FROM artists WHERE id = $1`

	if _, err = repo.db.ExecContext(ctx, query, artistId); err != nil {
		repo.log.WithError(err).Error("failed to query delete artist")
		return
	}

	return
}
