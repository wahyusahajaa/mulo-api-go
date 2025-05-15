package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type songRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewSongRepository(db *database.DB, log *logrus.Logger) contracts.SongRepository {
	return &songRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *songRepository) FindAll(ctx context.Context, pageSize int, offset int) (songs []models.Song, err error) {
	query := `
		SELECT 
			s.id,
			s.title,
			s.audio,
			s.duration,
			s.image,
			al.id as album_id ,
			al."name" as album_name,
			al.slug as album_slug,
			al.image as album_image,
			ar.id as artist_id,
			ar.name as artist_name,
			ar.slug as artist_slug,
			ar.image as artist_image
		FROM songs s  
		INNER JOIN albums al ON al.id = s.album_id
		INNER JOIN artists ar ON ar.id = al.artist_id
		ORDER BY s.id desc 
		LIMIT $1 OFFSET $2
	`
	args := []any{pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "FindAll", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		song := models.Song{}
		if err := rows.Scan(
			&song.Id,
			&song.Title,
			&song.Audio,
			&song.Duration,
			&song.Image,
			&song.Album.Id,
			&song.Album.Name,
			&song.Album.Slug,
			&song.Album.Image,
			&song.Album.Artist.Id,
			&song.Album.Artist.Name,
			&song.Album.Artist.Slug,
			&song.Album.Artist.Image,
		); err != nil {
			utils.LogError(repo.log, ctx, "song_repo", "FindAll", err)
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (repo *songRepository) FindCount(ctx context.Context) (total int, err error) {
	query := `SELECT COUNT(*) FROM songs`
	if err = repo.db.QueryRowContext(ctx, query).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "FindCount", err)
		return
	}
	return
}

func (repo *songRepository) FindSongById(ctx context.Context, id int) (song *models.Song, err error) {
	query := `
		SELECT 
			s.id,
			s.title,
			s.audio,
			s.duration,
			s.image,
			al.id as album_id ,
			al."name" as album_name,
			al.slug as album_slug,
			al.image as album_image,
			ar.id as artist_id,
			ar.name as artist_name,
			ar.slug as artist_slug,
			ar.image as artist_image
		FROM songs s  
		INNER JOIN albums al ON al.id = s.album_id
		INNER JOIN artists ar ON ar.id = al.artist_id
		WHERE s.id = $1
	`
	song = &models.Song{}

	if err = repo.db.QueryRowContext(ctx, query, id).Scan(
		&song.Id,
		&song.Title,
		&song.Audio,
		&song.Duration,
		&song.Image,
		&song.Album.Id,
		&song.Album.Name,
		&song.Album.Slug,
		&song.Album.Image,
		&song.Album.Artist.Id,
		&song.Album.Artist.Name,
		&song.Album.Artist.Slug,
		&song.Album.Artist.Image,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarn(repo.log, ctx, "song_repo", "GetSongById", errs.NewNotFoundError("Song", "id", id))
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "song_repo", "GetSongById", err)
		return nil, err
	}

	return song, nil
}

func (repo *songRepository) FindExistsSongById(ctx context.Context, id int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM songs WHERE id = $1)`
	if err = repo.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "FindExistsSongById", err)
		return
	}

	return
}

func (repo *songRepository) Store(ctx context.Context, input models.CreateSongInput) (err error) {
	query := `INSERT INTO songs(album_id, title, audio, duration, image) VALUES($1, $2, $3, $4, $5)`
	args := []any{input.AlbumId, input.Title, input.Audio, input.Duration, input.Image}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "Store", err)
		return
	}

	return
}

func (repo *songRepository) Update(ctx context.Context, input models.CreateSongInput, id int) (err error) {
	query := `UPDATE songs SET album_id = $1, audio = $2, title= $3, duration = $4, image = $5 WHERE id = $6`
	args := []any{input.AlbumId, input.Audio, input.Title, input.Duration, input.Image, id}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "Update", err)
		return
	}

	return
}

func (repo *songRepository) Delete(ctx context.Context, id int) (err error) {
	query := `DELETE FROM songs WHERE id = $1`

	if _, err = repo.db.ExecContext(ctx, query, id); err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "Delete", err)
		return
	}

	return
}

func (repo *songRepository) FindSongsByAlbumId(ctx context.Context, albumId int, pageSize int, offset int) (songs []models.Song, err error) {
	query := `
		SELECT 
			s.id,
			s.title,
			s.audio,
			s.duration,
			s.image,
			al.id as album_id ,
			al."name" as album_name,
			al.slug as album_slug,
			al.image as album_image,
			ar.id as artist_id,
			ar.name as artist_name,
			ar.slug as artist_slug,
			ar.image as artist_image
		FROM songs s  
		INNER JOIN albums al ON al.id = s.album_id
		INNER JOIN artists ar ON ar.id = al.artist_id
		WHERE 
			s.album_id = $1
		ORDER BY s.id desc 
		LIMIT $2 OFFSET $3
	`
	args := []any{albumId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "FindSongsByAlbumId", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		song := models.Song{}
		if err := rows.Scan(
			&song.Id,
			&song.Title,
			&song.Audio,
			&song.Duration,
			&song.Image,
			&song.Album.Id,
			&song.Album.Name,
			&song.Album.Slug,
			&song.Album.Image,
			&song.Album.Artist.Id,
			&song.Album.Artist.Name,
			&song.Album.Artist.Slug,
			&song.Album.Artist.Image,
		); err != nil {
			utils.LogError(repo.log, ctx, "song_repo", "FindSongsByAlbumId", err)
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (repo *songRepository) FindCountSongsByAlbumId(ctx context.Context, albumId int) (total int, err error) {
	query := `SELECT COUNT(*) FROM songs WHERE album_id = $1`
	if err = repo.db.QueryRowContext(ctx, query, albumId).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "song_repo", "FindCountSongsByAlbumId", err)
		return
	}
	return
}
