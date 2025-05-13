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

type albumRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewAlbumRepository(db *database.DB, log *logrus.Logger) contracts.AlbumRepository {
	return &albumRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *albumRepository) FindAll(ctx context.Context, pageSize int, offset int) (albums []models.Album, err error) {
	query := `SELECT id, artist_id, name, slug, image FROM albums ORDER BY id DESC LIMIT $1 OFFSET $2`
	args := []any{pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "FindAll", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		album := models.Album{}
		if err = rows.Scan(
			&album.Id,
			&album.ArtistId,
			&album.Name,
			&album.Slug,
			&album.Image,
		); err != nil {
			utils.LogError(repo.log, ctx, "album_repo", "FindAll", err)
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func (repo *albumRepository) FindAlbumById(ctx context.Context, id int) (album *models.AlbumWithArtist, err error) {
	query := `
		SELECT 
			al.id as album_id,
			al.name as album_name, 
			al.slug as album_slug,  
			al.image as album_image,
			ar.id as artist_id,
			ar.name as artist_name,
			ar.slug as artist_slug,
			ar.image as artist_image
		FROM 
			albums al 
		INNER JOIN artists ar ON ar.id = al.artist_id 
		WHERE al.id = $1
	`

	album = &models.AlbumWithArtist{}
	if err = repo.db.QueryRowContext(ctx, query, id).Scan(
		&album.Id,
		&album.Name,
		&album.Slug,
		&album.Image,
		&album.Artist.Id,
		&album.Artist.Name,
		&album.Artist.Slug,
		&album.Artist.Image,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFoundErr := utils.NotFoundError{Resource: "Album", Id: id}
			utils.LogWarn(repo.log, ctx, "album_repo", "FindAlbumById", notFoundErr)
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "album_repo", "FindAlbumById", err)
		return nil, err
	}

	return
}

func (repo *albumRepository) FindCount(ctx context.Context) (total int, err error) {
	query := `SELECT COUNT(*) FROM albums`

	if err = repo.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "FindCount", err)
		return
	}

	return
}

func (repo *albumRepository) FindExistsAlbumById(ctx context.Context, id int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM albums WHERE id = $1)`

	if err = repo.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "FindExistsAlbumById", err)
		return
	}

	return
}

func (repo *albumRepository) FindExistsAlbumBySlug(ctx context.Context, slug string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM albums WHERE slug = $1)`

	if err = repo.db.QueryRowContext(ctx, query, slug).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "FindExistsAlbumBySlug", err)
		return
	}

	return
}

func (repo *albumRepository) Store(ctx context.Context, input models.CreateAlbumInput) (err error) {
	query := `INSERT INTO albums(artist_id, name, slug, image) VALUES($1, $2, $3, $4)`
	args := []any{input.ArtistId, input.Name, input.Slug, input.Image}

	if _, err = repo.db.QueryContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "Store", err)
		return
	}

	return
}

func (repo *albumRepository) Update(ctx context.Context, input models.CreateAlbumInput, id int) (err error) {
	query := `UPDATE albums SET name = $1, artist_id = $2, slug = $3, image = $4 WHERE id = $5`
	args := []any{input.Name, input.ArtistId, input.Slug, input.Image, id}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "Update", err)
		return
	}

	return
}

func (repo *albumRepository) Delete(ctx context.Context, id int) (err error) {
	query := `DELETE FROM albums WHERE id = $1`

	if _, err = repo.db.ExecContext(ctx, query, id); err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "Delete", err)
		return
	}

	return
}

func (repo *albumRepository) FindAlbumsByArtistId(ctx context.Context, artistId int) (albums []models.Album, err error) {
	query := `SELECT id, artist_id, name, slug, image FROM albums WHERE artist_id = $1`

	rows, err := repo.db.QueryContext(ctx, query, artistId)
	if err != nil {
		utils.LogError(repo.log, ctx, "album_repo", "FindAlbumsByArtistId", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		album := models.Album{}
		if err := rows.Scan(&album.Id, &album.ArtistId, &album.Name, &album.Slug, &album.Image); err != nil {
			utils.LogError(repo.log, ctx, "album_repo", "FindAlbumsByArtistId", err)
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}
