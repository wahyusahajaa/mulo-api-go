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

func (repo *authRepository) Store(ctx context.Context, fullname, username, email, password, code string) (err error) {
	tx, err := repo.db.Begin()

	if err != nil {
		repo.log.WithError(err).Error("failed to begin transactiion")
		return errors.New("failed to begin transaction insert users")
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
		repo.log.WithError(err).Error("failed to query insert users")
		return errors.New("failed to query insert user")

	}

	if _, err = tx.ExecContext(ctx, queryInsertVerifyCode, userId, code); err != nil {
		repo.log.WithError(err).Error("failed to query insert user_verified")
		return errors.New("failed to query insert user_verified")
	}

	return nil
}

func (repo *authRepository) FindUserVerifiedByCode(ctx context.Context, code string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM user_verified WHERE code = $1)`

	if err := repo.db.QueryRowContext(ctx, query, code).Scan(&exists); err != nil {
		repo.log.WithError(err).Error("failed to query user_verified")
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) FindUserDuplicateEmail(ctx context.Context, email string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	if err := repo.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		repo.log.WithError(err).Error("failed to query user duplicate email")
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) FindUserDuplicateUsername(ctx context.Context, username string) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`

	if err := repo.db.QueryRowContext(ctx, query, username).Scan(&exists); err != nil {
		repo.log.WithError(err).Error("failed to query user duplicate username")
		return false, err
	}

	return exists, nil
}

func (repo *authRepository) FindUserByEmail(ctx context.Context, email string) (user *models.User, err error) {
	query := `SELECT id, full_name, email, username, password, role, image, email_verified_at FROM users WHERE email = $1`
	user = &models.User{}

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

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.log.WithField("email", email).Warn("user not found")
			return nil, nil
		}

		repo.log.WithError(err).Error("failed to query user")
		return nil, err
	}

	return user, nil
}

func (repo *authRepository) StoreUserVerifyCode(ctx context.Context, userId int, code string) (err error) {
	query := `INSERT INTO user_verified(user_id, code) VALUES($1, $2);`

	if _, err := repo.db.ExecContext(ctx, query, userId, code); err != nil {
		repo.log.WithError(err).Error("failed to query insert user_verified")
		return err
	}

	return nil
}

func (repo *authRepository) FindUserVerifiedByUserIdAndCode(ctx context.Context, userId int, code string) (userVerified *models.UserVerified, err error) {
	query := `SELECT uv.code, uv.expired_at FROM user_verified uv INNER JOIN users u ON u.id = uv.user_id WHERE uv.user_id = $1 AND uv.code = $2`

	userVerified = &models.UserVerified{}

	err = repo.db.QueryRowContext(ctx, query, userId, code).Scan(
		&userVerified.Code,
		&userVerified.ExpiredAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.log.WithFields(logrus.Fields{
				"user_id": userId,
				"code":    code,
			}).Warn("user verified not found")

			return nil, nil
		}

		repo.log.WithError(err).Error("failed to query user_verified")
		return nil, err
	}

	return userVerified, nil
}

func (repo *authRepository) UpdateUserVerifiedAt(ctx context.Context, userId int) (err error) {
	query := `UPDATE users SET email_verified_at = $1 WHERE id = $2`
	currentTime := time.Now()

	if _, err := repo.db.ExecContext(ctx, query, currentTime, userId); err != nil {
		repo.log.WithError(err).Error("failed to query update users")
		return err
	}

	return nil
}
