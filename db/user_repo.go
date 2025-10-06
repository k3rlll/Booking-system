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

func (r *UserRepository) NewUser(ctx context.Context, name string, email string) (User, error) {

	var userRes User

	if name == "" {
		return User{}, functions.ErrBadRequest
	}
	if !strings.Contains(email, "@") {
		return User{}, functions.ErrBadRequest
	}

	tag, err := r.pool.Exec(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return User{}, err
	}

	if tag.RowsAffected() == 0 {
		return User{}, functions.ErrUserAlreadyCreated

	}
	_ = r.pool.QueryRow(ctx,
		"SELECT user_id, name, email FROM users WHERE email=$1", email).Scan(&userRes.Id, &userRes.Name, &userRes.Email)

	return userRes, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (User, error) {

	var u User

	if id == 0 {
		return User{}, functions.ErrBadRequest
	}

	err := r.pool.QueryRow(ctx,
		"SELECT user_id, name, email FROM users WHERE user_id=$1", id).Scan(&u.Id, &u.Name, &u.Email)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *UserRepository) GetUserID(ctx context.Context, email string) (int, error) {
	var ID int
	err := r.pool.QueryRow(ctx,
		"select user_id from users where email=$1", email).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil

}
