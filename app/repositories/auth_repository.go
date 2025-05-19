package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	if _, err = repo.db.ExecContext(ctx, query, time.Now(), userId); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "UpdateUserVerifiedAt", err)
		return
	}

	return
}

func (repo *authRepository) IsRefreshTokenValid(ctx context.Context, userID int, token string) (valid bool, err error) {
	var revoked bool
	var expiredAt time.Time
	query := `SELECT revoked, expires_at FROM refresh_tokens WHERE user_id = $1 AND token = $2`

	if err = repo.db.QueryRowContext(ctx, query, userID, token).Scan(&revoked, &expiredAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nfErr := errs.NewNotFoundErrorWithMsg(fmt.Sprintf(`refresh_tokens with user_id: %d and token: %s does not exists.`, userID, token))
			utils.LogWarn(repo.log, ctx, "auth_repo", "FindExistsRefreshToken", nfErr)
			return false, nil
		}

		utils.LogWarn(repo.log, ctx, "auth_repo", "FindExistsRefreshToken", err)
		return false, err
	}

	// If has revoked or refresh token has expired set valid to false
	if revoked || expiredAt.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}

func (repo *authRepository) StoreRefreshToken(ctx context.Context, userID int, token string, expiredAt time.Time) (err error) {
	query := ` INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`

	if _, err = repo.db.ExecContext(ctx, query, userID, token, expiredAt); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindExistsRefreshToken", err)
		return
	}

	return
}

func (repo *authRepository) UpdateRefreshToken(ctx context.Context, userID int, token string) (err error) {
	query := `UPDATE refresh_tokens SET revoked = TRUE, revoked_at = NOW() WHERE user_id = $1 AND token = $2 AND revoked = FALSE`

	if _, err = repo.db.ExecContext(ctx, query, userID, token); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "UpdateRefreshToken", err)
		return
	}

	return
}

func (repo *authRepository) DeleteRefreshToken(ctx context.Context, token string) (err error) {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	if _, err = repo.db.ExecContext(ctx, query, token); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "DeleteRefreshToken", err)
		return
	}
	return
}

func (repo *authRepository) StoreUserWithOAuthAccount(ctx context.Context, input models.OAuthAccountInput) (userID int, err error) {
	tx, err := repo.db.Begin()
	if err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreOAuthAccounts", err)
		return 0, err
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

	userQuery := `INSERT INTO users(full_name, username, email, image, role, email_verified_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	userArgs := []any{input.Fullname, input.Username, input.Email, input.Image, "member", time.Now()}
	if err = tx.QueryRowContext(ctx, userQuery, userArgs...).Scan(&userID); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreOAuthAccounts", err)
		return 0, err
	}

	userAccountQuery := `INSERT INTO oauth_accounts(user_id, provider, provider_user_id) VALUES($1, $2, $3);`
	userAccountArgs := []any{userID, input.Provider, input.ID}
	if _, err = tx.ExecContext(ctx, userAccountQuery, userAccountArgs...); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreOAuthAccounts", err)
		return 0, err
	}

	return userID, nil
}

func (repo *authRepository) FindOAuthAccount(ctx context.Context, provider string, providerUserID string) (oAuthAccount *models.OAuthAccount, err error) {
	query := `SELECT id, user_id, provider, provider_user_id FROM oauth_accounts WHERE provider = $1 AND provider_user_id = $2`
	args := []any{provider, providerUserID}
	oAuthAccount = &models.OAuthAccount{}

	if err := repo.db.QueryRowContext(ctx, query, args...).Scan(
		&oAuthAccount.ID,
		&oAuthAccount.UserID,
		&oAuthAccount.Provider,
		&oAuthAccount.ProviderUserID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		utils.LogError(repo.log, ctx, "auth_repo", "FindOAuthAccount", err)
		return nil, err
	}

	return oAuthAccount, nil
}

func (repo *authRepository) StoreOAuthAccount(ctx context.Context, userID int, providerID, providerUserID string) (err error) {
	tx, err := repo.db.Begin()
	if err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreOAuthAccount", err)
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

	queryUpdateUser := `UPDATE users SET email_verified_at = $1 WHERE id = $2 AND email_verified_at IS NULL`
	if _, err = tx.ExecContext(ctx, queryUpdateUser, time.Now(), userID); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreOAuthAccount", err)
		return
	}

	userAccountQuery := `INSERT INTO oauth_accounts(user_id, provider, provider_user_id) VALUES($1, $2, $3);`
	userAccountArgs := []any{userID, providerID, providerUserID}
	if _, err = tx.ExecContext(ctx, userAccountQuery, userAccountArgs...); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "StoreOAuthAccount", err)
		return
	}

	return
}

func (repo *authRepository) FindExistsOauthAccount(ctx context.Context, userID int) (exists bool, err error) {
	query := `SELECT EXISTS (SELECT 1 FROM oauth_accounts WHERE user_id = $1)`
	if err = repo.db.QueryRowContext(ctx, query, userID).Scan(&exists); err != nil {
		utils.LogError(repo.log, ctx, "auth_repo", "FindExistsOauthAccount", err)
		return
	}

	return
}
