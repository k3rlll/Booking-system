package db

import (
	"context"
	"rest_api/functions"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) NewUser(ctx context.Context, name string, email string) error {

	if name == "" {
		return functions.ErrBadRequest
	}
	if !strings.Contains(email, "@") {
		return functions.ErrBadRequest
	}

	tag, err := r.pool.Exec(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return functions.ErrUserAlreadyCreated

	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (User, error) {

	var u User

	if id == 0 {
		return User{}, functions.ErrBadRequest
	}

	err := r.pool.QueryRow(ctx,
		"SELECT id, name, email FROM users WHERE id=$1", id).Scan(&u.Name, &u.Email)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (r *UserRepository) GetUserID(ctx context.Context, email string) (int, error) {
	var ID int
	err := r.pool.QueryRow(ctx,
		"select id from users where email=$1", email).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil

}
