package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Seats struct {
	Number int `json:"number"`
	Row    int `json:"row"`

	Is_reserved bool `json:"is_reserved"`
}

type SeatsRepository struct {
	pool *pgxpool.Pool
}

func NewSeatRepository(pool *pgxpool.Pool) *SeatsRepository {
	return &SeatsRepository{
		pool: pool,
	}
}

func (r *SeatsRepository) GetAllSeats(ctx context.Context) ([]Seats, error) {
	var s Seats
	rows, err := r.pool.Query(ctx,
		"SELECT number, row, is_reserved FROM seats ")
	if err != nil {
		return []Seats{}, err
	}

	defer rows.Close()

	var seats []Seats

	for rows.Next() {
		err = rows.Scan(&s.Number, &s.Row, &s.Is_reserved)
		if err != nil {
			return []Seats{}, err
		}
		seats = append(seats, s)

	}

	return seats, nil
}

func (r *SeatsRepository) GetFreeSeats(ctx context.Context) ([]Seats, error) {
	var s Seats
	rows, err := r.pool.Query(ctx,
		"SELECT number, row, is_reserved FROM seats WHERE is_reserved=$1", false)
	if err != nil {
		return []Seats{}, err
	}

	defer rows.Close()

	var seats []Seats

	for rows.Next() {
		err = rows.Scan(&s.Number, &s.Row, &s.Is_reserved)
		if err != nil {
			return []Seats{}, err
		}
		seats = append(seats, s)

	}

	return seats, nil

}

func (r *SeatsRepository) GetReservedSeats(ctx context.Context) ([]Seats, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT number, row, is_reserved FROM seats WHERE is_reserved=$1", true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []Seats

	for rows.Next() {
		var s Seats
		err = rows.Scan(&s.Number, &s.Row, &s.Is_reserved)
		if err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil

}
