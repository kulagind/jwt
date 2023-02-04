package repo

import (
	"context"
	"errors"
	"fmt"
	"jwt/internal/constants"
	"jwt/internal/models"
	"jwt/pkg/helpers/utils"
	"time"

	"github.com/google/uuid"
)

type userRepo struct{}

var userRepoSingleton *userRepo

func GetUserRepo() *userRepo {
	if userRepoSingleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if userRepoSingleton == nil {
			userRepoSingleton = &userRepo{}
		}
	}

	return userRepoSingleton
}

const createUserQuery = `
	insert into users 
	(id, email, password, name, tokenhash, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7)
`

func (r *userRepo) Create(
	ctx context.Context,
	candidate *models.User,
) (*models.UserResponse, error) {
	newId := uuid.NewString()
	newTokenHash := utils.GenerateRandomString(15)
	_, err := GetInstance().Db.Exec(ctx, createUserQuery,
		newId,
		candidate.Email,
		candidate.Password,
		"",
		newTokenHash,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		Id:    newId,
		Name:  candidate.Name,
		Email: candidate.Email,
	}, nil
}

func (r *userRepo) UpdateTokenhash(ctx context.Context, user *models.User) error {
	tokenhash := utils.GenerateRandomString(15)
	_, err := GetInstance().Db.Exec(ctx, `
		update users set tokenhash=$1 where id=$2;
	`, tokenhash, user.Id)
	if err != nil {
		return err
	}
	user.TokenHash = tokenhash
	return nil
}

func (r *userRepo) FindBy(ctx context.Context, field string, value string) (*models.UserResponse, error) {
	query := GetQueryFindUserBy(field)
	user := &models.UserResponse{}
	rows, err := GetInstance().Db.Query(ctx, query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Name)
		if err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil, errors.New(constants.NO_ROWS)
	}

	return user, nil
}

func GetQueryFindUserBy(field string) string {
	return fmt.Sprintf(
		"select id, email, name from users where %s = $1 limit 1;",
		field,
	)
}

func (r *userRepo) PrivateFindBy(ctx context.Context, field string, value string) (*models.User, error) {
	query := GetQueryPrivateFindUserBy(field)
	user := &models.User{}
	rows, err := GetInstance().Db.Query(ctx, query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.TokenHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil, errors.New(constants.NO_ROWS)
	}

	return user, nil
}

func GetQueryPrivateFindUserBy(field string) string {
	return fmt.Sprintf(
		"select id, email, name, password, tokenhash, created_at, updated_at from users where %s = $1 limit 1;",
		field,
	)
}
