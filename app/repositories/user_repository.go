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

func (repo *userRepository) FindExistsUserByUserID(ctx context.Context, userID int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`

	if err = repo.db.QueryRowContext(ctx, query, userID).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "FindExistsUserById", err)
		return
	}

	return
}

func (repo *userRepository) FindUserVerifiedByCode(ctx context.Context, code string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM user_verified WHERE code = $1)`

	if err := repo.db.QueryRowContext(ctx, query, code).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindUserVerifiedByCode", err)
		return false, err
	}

	return exists, nil
}

func (repo *userRepository) FindUserExistsByEmail(ctx context.Context, email string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	if err = repo.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindUserExistsByEmail", err)
		return
	}

	return
}

func (repo *userRepository) FindUserExistsByUsername(ctx context.Context, username string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	if err = repo.db.QueryRowContext(ctx, query, username).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindUserExistsByUsername", err)
		return
	}

	return
}

func (repo *userRepository) FindUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	query := `SELECT id, full_name, email, username, password, role, image, email_verified_at FROM users WHERE email = $1`

	user = &models.User{}
	if err = repo.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Image,
		&user.EmailVerifiedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarn(repo.log, ctx, "auth_repo", "FindUserByEmail", errs.NewNotFoundError("User", "email", email))
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "auth_repo", "FindUserByEmail", err)
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) FindUserByUserID(ctx context.Context, userID int) (user *models.User, err error) {
	query := `SELECT id, full_name, email, username, password, role, image, email_verified_at FROM users WHERE id = $1`

	user = &models.User{}
	if err = repo.db.QueryRowContext(ctx, query, userID).Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Image,
		&user.EmailVerifiedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarn(repo.log, ctx, "auth_repo", "FindUserByUserID", errs.NewNotFoundError("User", "id", userID))
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "auth_repo", "FindUserByUserID", err)
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) FindUserVerifiedByUserIDAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error) {
	query := `SELECT uv.code, uv.expired_at FROM user_verified uv INNER JOIN users u ON u.id = uv.user_id WHERE uv.user_id = $1 AND uv.code = $2 AND u.email_verified_at IS NULL`

	userVerified = &models.UserVerified{}
	if err = repo.db.QueryRowContext(ctx, query, userId, code).Scan(&userVerified.Code, &userVerified.ExpiredAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarn(repo.log, ctx, "auth_repo", "FindUserVerifiedByUserIdAndCode", errs.NewNotFoundError("UserVerified", "code", code))
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "auth_repo", "FindUserVerifiedByUserIdAndCode", err)
		return nil, err
	}

	return userVerified, nil
}

func (repo *userRepository) Count(ctx context.Context) (total int, err error) {
	query := `SELECT COUNT(*) FROM users WHERE role != $1`

	if err = repo.db.QueryRowContext(ctx, query, "admin").Scan(&total); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "Count", err)
		return
	}

	return
}

func (repo *userRepository) Update(ctx context.Context, input models.CreateUserInput, userID int) (err error) {
	query := `UPDATE users SET full_name = $1, image = $2 WHERE id = $3`
	args := []any{input.Fullname, input.Image, userID}

	if _, err := repo.db.ExecContext(ctx, query, args...); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "Update", err)
		return err
	}

	return nil
}

func (repo *userRepository) Delete(ctx context.Context, userID int) (err error) {
	query := `DELETE FROM users WHERE id = $1`

	if _, err := repo.db.ExecContext(ctx, query, userID); err != nil {
		utils.LogError(repo.log, ctx, "user_repo", "Delete", err)
		return err
	}

	return nil
}
