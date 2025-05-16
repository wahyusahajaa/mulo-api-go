package repositories

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type favoriteRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewFavoriteRepository(db *database.DB, log *logrus.Logger) contracts.FavoriteRepository {
	return &favoriteRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *favoriteRepository) FindFavoriteSongsByUserID(ctx context.Context, userId, pageSize, offset int) (songs []models.Song, err error) {
	query := `
		SELECT 
			s.id AS song_id,
			s.album_id AS song_album_id,
			s.title AS song_title,
			s.audio AS song_audio,
			s.duration AS song_duration,
			s.image AS song_image,
			al.id AS album_id,
			al.artist_id AS album_artist_id,
			al.name AS album_name,
			al.slug AS album_slug,
			al.image AS album_image,
			ar.id AS artist_id,
			ar.name AS artist_name,
			ar.slug AS artist_slug,
			ar.image AS artist_image
		FROM song_favorites sf 
		INNER JOIN songs s on s.id  = sf.song_id
		INNER JOIN albums al on al.id = s.album_id 
		INNER JOIN artists ar on ar.id  = al.artist_id 
		WHERE
			user_id = $1
		ORDER by sf.created_at desc
		LIMIT $2 offset $3
	`
	args := []any{userId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "favorite_repo", "FindExistsFavoriteSongBySongId", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var song models.Song
		if err = rows.Scan(
			&song.Id,
			&song.AlbumId,
			&song.Title,
			&song.Audio,
			&song.Duration,
			&song.Image,
			&song.Album.Id,
			&song.Album.ArtistId,
			&song.Album.Name,
			&song.Album.Slug,
			&song.Album.Image,
			&song.Album.Artist.Id,
			&song.Album.Artist.Name,
			&song.Album.Artist.Slug,
			&song.Album.Artist.Image,
		); err != nil {
			utils.LogError(repo.log, ctx, "favorite_repo", "FindExistsFavoriteSongBySongId", err)
			return
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (repo *favoriteRepository) FindCountFavoriteSongsByUserID(ctx context.Context, userId int) (total int, err error) {
	query := `SELECT COUNT(*) FROM song_favorites WHERE user_id = $1`
	if err = repo.db.QueryRowContext(ctx, query, userId).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "favorite_repo", "FindCountFavoriteSongsByUserId", err)
		return
	}
	return
}

func (repo *favoriteRepository) FindExistsFavoriteSongBySongID(ctx context.Context, userId int, songId int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM song_favorites WHERE user_id = $1 AND song_id = $2)`
	args := []any{userId, songId}

	if err = repo.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "favorite_repo", "favorite_service", err)
		return
	}

	return
}

func (repo *favoriteRepository) StoreFavoriteSong(ctx context.Context, userId int, songId int) (err error) {
	query := `INSERT INTO song_favorites(user_id, song_id) VALUES($1, $2)`
	args := []any{userId, songId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "favorite_repo", "Store", err)
		return
	}

	return
}

func (repo *favoriteRepository) DeleteFavoriteSong(ctx context.Context, userId int, songId int) (err error) {
	query := `DELETE FROM song_favorites WHERE user_id = $1 AND song_id = $2`
	args := []any{userId, songId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "favorite_repo", "DeleteFavoriteSong", err)
		return
	}

	return
}
