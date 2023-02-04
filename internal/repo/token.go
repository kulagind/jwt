package repo

import (
	"context"
	"errors"
	"jwt/internal/models"
	"jwt/pkg/helpers/pg"
	"time"

	"github.com/google/uuid"
)

type tokenRepo struct{}

var tokenRepoSingleton *tokenRepo

func GetTokenRepo() *tokenRepo {
	if tokenRepoSingleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if tokenRepoSingleton == nil {
			tokenRepoSingleton = &tokenRepo{}
		}
	}

	return tokenRepoSingleton
}

func (r *tokenRepo) UpdateRefresh(
	ctx context.Context,
	oldToken string,
	newToken string,
) error {
	newId := uuid.NewString()
	_, err := GetInstance().Db.Exec(ctx, `
			insert into updated_tokens 
			(id, old_token, new_token, created_at, updated_at) 
			values ($1, $2, $3, $4, $5)
		`,
		newId,
		oldToken,
		newToken,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *tokenRepo) CheckRefresh(ctx context.Context, token string) error {
	repoUpdated := &models.RepoUpdatedTokens{}
	rows, err := GetInstance().Db.Query(ctx, `
		select old_token, new_token from updated_tokens where old_token=$1 limit 1;
	`, token)
	if err != nil && !pg.CheckSqlError(err, "no rows in result set") {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&repoUpdated.OldToken, &repoUpdated.NewToken)
		if err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if repoUpdated.OldToken != "" {
		err = GetTokenRepo().BlockRefreshToken(ctx, repoUpdated.OldToken)
		if err != nil {
			return err
		}
		err = GetTokenRepo().BlockRefreshToken(ctx, repoUpdated.NewToken)
		if err != nil {
			return err
		}
		return errors.New("token was block as suspicious")
	}

	repoBlocked := &models.RepoBlocked{}
	rows, err = GetInstance().Db.Query(ctx, `
		select token from black_list where token=$1 limit 1;
	`, token)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&repoBlocked.Token)
		if err != nil {
			return err
		}
	}
	if repoBlocked.Token != "" {
		return errors.New("token was block as suspicious")
	}

	return nil
}

func (r *tokenRepo) BlockRefreshToken(
	ctx context.Context,
	token string,
) error {
	newId := uuid.NewString()
	_, err := GetInstance().Db.Exec(ctx, `
			insert into black_list 
			(id, token, created_at, updated_at) 
			values ($1, $2, $3, $4, $5)
		`,
		newId,
		token,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
