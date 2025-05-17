package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
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

func (repo *authRepository) StoreUserVerifyCode(ctx context.Context, userId int, code string) (err error) {
	query := `INSERT INTO user_verified(user_id, code) VALUES($1, $2);`

	if _, err = repo.db.ExecContext(ctx, query, userId, code); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreUserVerifyCode", err)
		return
	}

	return
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
