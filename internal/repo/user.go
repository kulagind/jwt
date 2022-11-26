package repo

import (
	"context"
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

func (r *userRepo) Create(
	ctx context.Context,
	candidate *models.User,
) (*models.UserResponse, error) {
	newId := uuid.NewString()
	_, err := getInstance().Db.Exec(ctx, `
			insert into users 
			(id, email, password, name, tokenhash, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7)
		`,
		newId,
		candidate.Email,
		candidate.Password,
		candidate.Name,
		utils.GenerateRandomString(15),
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

func (r *userRepo) FindById(ctx context.Context, id string) (*models.UserResponse, error) {
	user := &models.UserResponse{}
	err := getInstance().Db.QueryRow(ctx, `
		select id, email, name from users where id=$1 limit 1;
	`, id).Scan(&user.Id, &user.Email, &user.Name)
	if err != nil {
		return nil, err
	}
	return user, nil
}
