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

type playlistRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewPlaylistRepository(db *database.DB, log *logrus.Logger) contracts.PlaylistRepository {
	return &playlistRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *playlistRepository) FindAll(ctx context.Context, role string, userId, pageSize, offset int) (playlists []models.Playlist, err error) {
	query := `SELECT id, name FROM playlists`
	sort := ` ORDER BY id DESC`
	var args []any

	if role == "member" {
		query += ` WHERE user_id = $1` + ` LIMIT $2 OFFSET $3`
		args = []any{userId, pageSize, offset}
	} else {
		query += sort + ` LIMIT $1 OFFSET $2`
		args = []any{pageSize, offset}
	}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "FindAll", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		playlist := models.Playlist{}
		if err := rows.Scan(&playlist.Id, &playlist.Name); err != nil {
			utils.LogError(repo.log, ctx, "playlist_repo", "FindAll", err)
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (repo *playlistRepository) FindById(ctx context.Context, role string, userId, id int) (playlist *models.Playlist, err error) {
	query := `SELECT id, name FROM playlists WHERE id = $1`
	var args []any

	if role == "member" {
		query += ` AND user_id = $2`
		args = []any{id, userId}
	} else {
		args = []any{id}
	}

	playlist = &models.Playlist{}
	if err = repo.db.QueryRowContext(ctx, query, args...).Scan(&playlist.Id, &playlist.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarn(repo.log, ctx, "playlist_repo", "FindById", utils.NotFoundError{Resource: "Playlist", Id: id})
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "playlist_repo", "FindById", err)
		return nil, err
	}

	return playlist, nil
}

func (repo *playlistRepository) FindCount(ctx context.Context, role string, userId int) (total int, err error) {
	query := `SELECT COUNT(*) FROM playlists`
	var args []any

	if role == "member" {
		query += ` WHERE user_id = $1`
		args = []any{userId}
	}

	if err = repo.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "FindCount", err)
		return
	}

	return
}

func (repo *playlistRepository) FindExistsPlaylistById(ctx context.Context, userRole string, userId, playlistId int) (exists bool, err error) {
	condition := `WHERE id = $1`
	var args []any

	if userRole == "member" {
		condition += ` AND user_id = $2`
		args = []any{playlistId, userId}
	} else {
		args = []any{playlistId}
	}

	query := `SELECT EXISTS (SELECT 1 FROM playlists ` + condition + `)`

	if err = repo.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "FindExistsPlaylistById", err)
		return
	}

	return
}

func (repo *playlistRepository) Store(ctx context.Context, input models.CreatePlaylistInput) (err error) {
	query := `INSERT INTO playlists(user_id, name) VALUES($1, $2)`
	args := []any{input.UserId, input.Name}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "Store", err)
		return
	}

	return
}

func (repo *playlistRepository) Update(ctx context.Context, input models.CreatePlaylistInput, id int) (err error) {
	query := `UPDATE playlists SET name = $1 WHERE user_id = $2 AND id = $3`
	args := []any{input.Name, input.UserId, id}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "Update", err)
		return
	}

	return
}

func (repo *playlistRepository) Delete(ctx context.Context, userRole string, userId, playlistId int) (err error) {
	query := `DELETE FROM playlists WHERE id = $1`
	var args []any

	if userRole == "member" {
		query += ` AND user_id = $2`
		args = []any{playlistId, userId}
	} else {
		args = []any{playlistId}
	}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "Delete", err)
		return
	}

	return
}

func (repo *playlistRepository) FindPlaylistSongs(ctx context.Context, playlistId int, pageSize int, offset int) (songs []models.Song, err error) {
	query := `
		SELECT 
			s.id,
			s.title,
			s.audio,
			s.duration,
			s.image,
			al.id as album_id ,
			al.name as album_name,
			al.slug as album_slug,
			al.image as album_image,
			ar.id as artist_id,
			ar.name as artist_name,
			ar.slug as artist_slug,
			ar.image as artist_image
		FROM playlist_songs ps
		INNER JOIN songs s on s.id  = ps.song_id
		INNER JOIN albums al ON al.id = s.album_id
		INNER JOIN artists ar ON ar.id = al.artist_id
		WHERE ps.playlist_id = $1
		ORDER BY s.id desc 
		LIMIT $2 OFFSET $3
	`
	args := []any{playlistId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "FindPlaylistSongs", err)
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
			utils.LogError(repo.log, ctx, "playlist_repo", "FindPlaylistSongs", err)
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (repo *playlistRepository) StorePlaylistSong(ctx context.Context, playlistId int, songId int) (err error) {
	query := `INSERT INTO playlist_songs(playlist_id, song_id) VALUES($1, $2)`
	args := []any{playlistId, songId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "StorePlaylistSong", err)
		return
	}

	return
}

func (repo *playlistRepository) FindExistsPlaylistSong(ctx context.Context, playlistId int, songId int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2)`
	args := []any{playlistId, songId}

	if err = repo.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "FindExistsPlaylistSongByPlaylistId", err)
		return
	}

	return
}

func (repo *playlistRepository) DeletePlaylistSong(ctx context.Context, playlistId int, songId int) (err error) {
	query := `DELETE FROM playlist_songs WHERE playlist_id = $1 AND song_id = $2`
	args := []any{playlistId, songId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "DeletePlaylistSong", err)
		return
	}

	return
}
