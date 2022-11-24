package repo

import (
	"context"
	"jwt/internal/models"
	"jwt/pkg/helpers/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func (repo *UserRepo) Create(
	ctx context.Context,
	name string,
	email string,
	password string,
) error {
	_, err := repo.db.Exec(ctx, `
			insert into users 
			(id, email, password, name, tokenhash, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7)
		`,
		uuid.NewString(),
		email,
		password,
		name,
		utils.GenerateRandomString(15),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) FindById(ctx context.Context, id string) (*models.UserResponse, error) {
	user := &models.UserResponse{}
	err := repo.db.QueryRow(ctx, `
		select id, email, name from users where id=$1 limit 1;
	`, id).Scan(&user.Id, &user.Email, &user.Name)
	if err != nil {
		return nil, err
	}
	return user, nil
}
