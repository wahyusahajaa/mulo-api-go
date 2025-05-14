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

type userRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewUserRepository(db *database.DB, log *logrus.Logger) contracts.UserRepository {
	return &userRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *userRepository) FindAll(ctx context.Context, pageSize, offset int) (users []models.User, err error) {
	query := `SELECT id, full_name, username, email, password, image, role, email_verified_at FROM users WHERE role != $1 ORDER BY id DESC LIMIT $2 OFFSET $3`
	rows, err := repo.db.QueryContext(ctx, query, "admin", pageSize, offset)

	if err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "FindAll", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(
			&user.Id,
			&user.Fullname,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Image,
			&user.Role,
			&user.EmailVerifiedAt,
		); err != nil {
			utils.LogError(repo.log, ctx, "user_repo", "FindAll", err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo *userRepository) FindExistsUserById(ctx context.Context, userId int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

	if err = repo.db.QueryRowContext(ctx, query, userId).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "FindExistsUserById", err)
		return
	}

	return
}

func (repo *userRepository) Count(ctx context.Context) (total int, err error) {
	query := `SELECT COUNT(*) FROM users WHERE role != $1`

	if err = repo.db.QueryRowContext(ctx, query, "admin").Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "Count", err)
		return
	}

	return
}

func (repo *userRepository) FindUserById(ctx context.Context, userId int) (user *models.User, err error) {
	query := `SELECT id, full_name, username, email, password, image, role, email_verified_at FROM users WHERE id = $1`
	user = &models.User{}
	if err = repo.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.Fullname,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Image,
		&user.Role,
		&user.EmailVerifiedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarn(repo.log, ctx, "user_repo", "FindUserById", errs.NewNotFoundError("User", "id", userId))
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "user_repo", "FindUserById", err)
		return nil, err
	}

	return
}

func (repo *userRepository) Update(ctx context.Context, input models.CreateUserInput, userId int) (err error) {
	query := `UPDATE users SET full_name = $1, image = $2 WHERE id = $3`
	args := []any{input.Fullname, input.Image, userId}

	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "Update", err)
		return err
	}

	return nil
}

func (repo *userRepository) Delete(ctx context.Context, userId int) (err error) {
	query := `DELETE FROM users WHERE id = $1`

	if _, err := repo.db.ExecContext(ctx, query, userId); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "Delete", err)
		return err
	}

	return nil
}
