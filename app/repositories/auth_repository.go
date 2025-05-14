package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type authRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewAuthRepository(db *database.DB, log *logrus.Logger) contracts.AuthRepository {
	return &authRepository{
		db:  db.DB,
		log: log,
	}
}

func (repo *authRepository) Store(ctx context.Context, input models.RegisterInput) (err error) {
	tx, err := repo.db.Begin()
	if err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "Store", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var userId int
	userQuery := `INSERT INTO users(full_name, username, email, password, role) VALUES($1, $2, $3, $4, $5) RETURNING id`
	userArgs := []any{input.Fullname, input.Username, input.Email, input.Password, "member"}
	if err = tx.QueryRowContext(ctx, userQuery, userArgs...).Scan(&userId); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "Store", err)
		return err
	}

	userVerifiedQuery := `INSERT INTO user_verified(user_id, code) VALUES($1, $2);`
	verifyArgs := []any{userId, input.Code}
	if _, err = tx.ExecContext(ctx, userVerifiedQuery, verifyArgs...); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "Store", err)
		return err
	}

	return nil
}

func (repo *authRepository) FindUserVerifiedByCode(ctx context.Context, code string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM user_verified WHERE code = $1)`

	if err := repo.db.QueryRowContext(ctx, query, code).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindUserVerifiedByCode", err)
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) FindUserExistsByEmail(ctx context.Context, email string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	if err = repo.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindUserExistsByEmail", err)
		return
	}

	return
}

func (repo *authRepository) FindUserExistsByUsername(ctx context.Context, username string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	if err = repo.db.QueryRowContext(ctx, query, username).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindUserExistsByUsername", err)
		return
	}

	return
}

func (repo *authRepository) FindUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
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

func (repo *authRepository) StoreUserVerifyCode(ctx context.Context, userId int, code string) (err error) {
	query := `INSERT INTO user_verified(user_id, code) VALUES($1, $2);`

	if _, err = repo.db.ExecContext(ctx, query, userId, code); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreUserVerifyCode", err)
		return
	}

	return
}

func (repo *authRepository) FindUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error) {
	query := `SELECT uv.code, uv.expired_at FROM user_verified uv INNER JOIN users u ON u.id = uv.user_id WHERE uv.user_id = $1 AND uv.code = $2`

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

func (repo *authRepository) UpdateUserVerifiedAt(ctx context.Context, userId int) (err error) {
	query := `UPDATE users SET email_verified_at = $1 WHERE id = $2`
	currentTime := time.Now()

	if _, err = repo.db.ExecContext(ctx, query, currentTime, userId); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "UpdateUserVerifiedAt", err)
		return
	}

	return
}
