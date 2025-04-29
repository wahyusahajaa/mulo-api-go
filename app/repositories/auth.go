package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/database"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *database.DB) contracts.AuthRepository {
	return &authRepository{
		db: db.DB,
	}
}

func (repo *authRepository) Store(ctx context.Context, fullname, username, email, password, code string) (err error) {
	tx, err := repo.db.Begin()

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
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
	queryInsertUser := `INSERT INTO users(full_name, username, email, password, role) VALUES($1, $2, $3, $4, $5) RETURNING id`
	queryInsertVerifyCode := `INSERT INTO user_verified(user_id, code) VALUES($1, $2);`

	if err = tx.QueryRowContext(ctx, queryInsertUser, fullname, username, email, password, "member").Scan(&userId); err != nil {
		return fmt.Errorf("failed while create user: %w", err)
	}

	if _, err = tx.ExecContext(ctx, queryInsertVerifyCode, userId, code); err != nil {
		return fmt.Errorf("failed while insert user verify code: %w", err)
	}

	return nil
}

func (repo *authRepository) FindUserVerifiedByCode(ctx context.Context, code string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM user_verified WHERE code = $1)`

	if err := repo.db.QueryRowContext(ctx, query, code).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) FindUserDuplicateEmail(ctx context.Context, email string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	if err := repo.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) FindUserDuplicateUsername(ctx context.Context, username string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	if err := repo.db.QueryRowContext(ctx, query, username).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) FindUserByEmail(ctx context.Context, email string) (user models.User, err error) {
	query := `SELECT id, full_name, email, username, password, role, image, email_verified_at FROM users WHERE email = $1`

	err = repo.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.Image,
		&user.EmailVerifiedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return user, nil
	}

	return
}

func (repo *authRepository) StoreVerifyCode(ctx context.Context, userId int, code string) (err error) {
	query := `INSERT INTO user_verified(user_id, code) VALUES($1, $2);`

	if _, err := repo.db.ExecContext(ctx, query, userId, code); err != nil {
		return err
	}

	return nil
}

func (repo *authRepository) FindUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users u INNER JOIN user_verified uv ON uv.user_id = u.id WHERE u.id = $1 AND uv.code = $2)`

	if err := repo.db.QueryRowContext(ctx, query, userId, code).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) UpdateUserVerifiedAt(ctx context.Context, userId int) (err error) {
	query := `UPDATE users SET email_verified_at = $1 WHERE id = $2`
	currentTime := time.Now()

	if _, err := repo.db.ExecContext(ctx, query, currentTime, userId); err != nil {
		return err
	}

	return nil
}
