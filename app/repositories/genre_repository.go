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
			utils.LogWarn(repo.log, ctx, "genre_repo", "FindGenreById", errs.NewNotFoundError("Genre", "id", id))
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
		utils.LogError(repo.log, ctx, "genre_repo", "FindExistsArtistGenreByGenreId", err)
		return
	}

	return
}

func (repo *genreRepository) StoreArtistGenre(ctx context.Context, artistId int, genreId int) (err error) {
	query := `INSERT INTO artist_genres(artist_id, genre_id) VALUES($1, $2)`
	args := []any{artistId, genreId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "StoreArtistGenre", err)
		return
	}

	return
}

func (repo *genreRepository) FindArtistGenres(ctx context.Context, artistId, pageSize, offset int) (genres []models.Genre, err error) {
	query := `SELECT g.id AS genre_id, g.name AS genre_name, g.image AS genre_image FROM artist_genres ag INNER JOIN genres g ON g.id = ag.genre_id WHERE ag.artist_id = $1 ORDER BY ag.created_at DESC LIMIT $2 OFFSET $3`
	args := []any{artistId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindArtistGenres", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.Id, &genre.Name, &genre.Image); err != nil {
			utils.LogError(repo.log, ctx, "genre_repo", "FindArtistGenres", err)
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
		utils.LogError(repo.log, ctx, "genre_repo", "DeleteArtistGenre", err)
		return err
	}

	return
}

func (repo *genreRepository) StoreSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	query := `INSERT INTO song_genres(song_id, genre_id) VALUES($1, $2)`
	args := []any{songId, genreId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "StoreSongGenre", err)
		return
	}

	return
}

func (repo *genreRepository) FindSongGenres(ctx context.Context, songId int, pageSize int, offset int) (genres []models.Genre, err error) {
	query := `SELECT g.id AS genre_id, g.name AS genre_name, g.image AS genre_image FROM song_genres sg INNER JOIN genres g ON g.id = sg.genre_id WHERE sg.song_id = $1 ORDER BY sg.created_at DESC LIMIT $2 OFFSET $3`
	args := []any{songId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindSongGenres", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.Id, &genre.Name, &genre.Image); err != nil {
			utils.LogError(repo.log, ctx, "genre_repo", "FindSongGenres", err)
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (repo *genreRepository) FindExistsSongGenreByGenreId(ctx context.Context, songId int, genreId int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM song_genres WHERE song_id = $1 AND genre_id = $2)`
	args := []any{songId, genreId}

	if err = repo.db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindExistsSongGenreByGenreId", err)
		return
	}

	return
}

func (repo *genreRepository) DeleteSongGenre(ctx context.Context, songId int, genreId int) (err error) {
	query := `DELETE FROM song_genres WHERE song_id = $1 AND genre_id = $2`
	args := []any{songId, genreId}

	if _, err = repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "DeleteArtistGenre", err)
		return err
	}

	return
}

func (repo *genreRepository) FindAllArtists(ctx context.Context, genreId int, pageSize int, offset int) (artists []models.Artist, err error) {
	query := `SELECT ar.id, ar.name, ar.slug, ar.image FROM artist_genres ag INNER JOIN artists ar ON ar.id = ag.artist_id WHERE ag.genre_id = $1 ORDER BY ag.created_at DESC LIMIT $2 OFFSET $3`
	args := []any{genreId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindAllArtists", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var artist models.Artist
		if err = rows.Scan(&artist.Id, &artist.Name, &artist.Slug, &artist.Image); err != nil {
			utils.LogError(repo.log, ctx, "genre_repo", "FindAllArtists", err)
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func (repo *genreRepository) FindCountArtists(ctx context.Context, genreId int) (total int, err error) {
	query := `SELECT COUNT(*) FROM artist_genres WHERE genre_id = $1`

	if err = repo.db.QueryRowContext(ctx, query, genreId).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindCountArtists", err)
		return total, err
	}

	return
}

func (repo *genreRepository) FindAllSongs(ctx context.Context, genreId int, pageSize int, offset int) (songs []models.Song, err error) {
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
		FROM song_genres sg
		INNER JOIN songs s ON s.id = sg.song_id
		INNER JOIN albums al ON al.id = s.album_id
		INNER JOIN artists ar ON ar.id = al.artist_id
		WHERE sg.genre_id = $1
		ORDER BY sg.created_at desc 
		LIMIT $2 OFFSET $3
	`
	args := []any{genreId, pageSize, offset}

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindAllSongs", err)
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
			utils.LogError(repo.log, ctx, "genre_repo", "FindAllSongs", err)
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (repo *genreRepository) FindCountSongs(ctx context.Context, genreId int) (total int, err error) {
	query := `SELECT COUNT(*) FROM song_genres WHERE genre_id = $1`

	if err = repo.db.QueryRowContext(ctx, query, genreId).Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "genre_repo", "FindCountSongs", err)
		return total, err
	}

	return
}
