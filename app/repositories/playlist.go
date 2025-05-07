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

func isValidRoles(role string) bool {
	roles := map[string]struct{}{
		"member": {},
		"admin":  {},
	}

	if _, ok := roles[role]; ok {
		return true
	}

	return false
}

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

func (repo *playlistRepository) FindAll(ctx context.Context, pageSize int, offset int) (playlists []models.Playlist, err error) {
	role := utils.GetRole(ctx)
	userId := utils.GetUserId(ctx)
	query := `SELECT id, name FROM playlists`
	sort := ` ORDER BY id DESC`
	var args []any

	if !isValidRoles(role) {
		err = fmt.Errorf("role not found while to query playlists")
		utils.LogError(repo.log, ctx, "playlist_repo", "FindAll", err)
		return nil, err
	}

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

func (repo *playlistRepository) FindById(ctx context.Context, id int) (playlist *models.Playlist, err error) {
	role := utils.GetRole(ctx)
	userId := utils.GetUserId(ctx)
	query := `SELECT id, name FROM playlists WHERE id = $1`
	var args []any

	if !isValidRoles(role) {
		err = fmt.Errorf("role not found while to query playlist")
		utils.LogError(repo.log, ctx, "playlist_repo", "FindById", err)
		return nil, err
	}

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

func (repo *playlistRepository) FindCount(ctx context.Context) (total int, err error) {
	role := utils.GetRole(ctx)
	userId := utils.GetUserId(ctx)
	query := `SELECT COUNT(*) FROM playlists`
	var args []any

	if !isValidRoles(role) {
		err = fmt.Errorf("role not found while to query count playlitst")
		utils.LogError(repo.log, ctx, "playlist_repo", "FindCount", err)
		return
	}

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

func (repo *playlistRepository) FindExistsPlaylistById(ctx context.Context, id int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM playlists WHERE id = $1)`

	if err = repo.db.QueryRowContext(ctx, query, id).Scan(&exists); err != nil {
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

func (repo *playlistRepository) Delete(ctx context.Context, id int) (err error) {
	query := `DELETE FROM playlists WHERE id = $1`

	if _, err = repo.db.ExecContext(ctx, query, id); err != nil {
		utils.LogError(repo.log, ctx, "playlist_repo", "Delete", err)
		return
	}

	return
}
