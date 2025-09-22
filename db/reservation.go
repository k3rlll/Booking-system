package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepositoty struct {
	pool *pgxpool.Pool
}

func (r *ReservationRepositoty) Reserve(pool *pgxpool.Pool) error {

	//tag, err := r.pool.Exec(context.Background(),
	//	"INSERT INTO ")

	return nil
}
