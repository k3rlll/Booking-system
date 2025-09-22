package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	id    int
	name  string
	email string
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) NewUser(pgx *pgxpool.Pool, name string, email string) error {

	tag, err := r.pool.Exec(context.Background(),
		"INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("failed to create a new user")

	}

	return nil
}

func (r *UserRepository) GetUserByID(pool *pgxpool.Pool, id int) (User, error) {

	var u User

	err := r.pool.QueryRow(context.Background(),
		"SELECT id, name, email FROM users WHERE id=$1", id).Scan(&u.name, &u.email)
	if err != nil {
		return User{}, err
	}

	return u, nil
}


