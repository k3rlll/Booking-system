package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Seats struct {
	Id     int
	Number int
	Row    int

	Is_reserved bool
}

type SeatsRepository struct {
	pool *pgxpool.Pool
}

func (r *SeatsRepository) GetAllSeats() ([]Seats, error) {
	var s Seats
	rows, err := r.pool.Query(context.Background(),
		"SELECT id, number, row, is_reserved FROM seats ")
	if err != nil {
		return []Seats{}, err
	}

	defer rows.Close()

	var seats []Seats

	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Number, &s.Row, &s.Is_reserved)
		if err != nil {
			return []Seats{}, err
		}
		seats = append(seats, s)

	}

	return seats, nil
}

func (r *SeatsRepository) GetFreeSeats() ([]Seats, error) {
	var s Seats
	rows, err := r.pool.Query(context.Background(),
		"SELECT id, number, row, is_reserved FROM seats WHERE is_reserved=$1", false)
	if err != nil {
		return []Seats{}, err
	}

	defer rows.Close()

	var seats []Seats

	for rows.Next() {
		err = rows.Scan(&s.Id, &s.Number, &s.Row, &s.Is_reserved)
		if err != nil {
			return []Seats{}, err
		}
		seats = append(seats, s)

	}

	return seats, nil

}

func (r *SeatsRepository) GetReservedSeats() ([]Seats, error) {
	rows, err := r.pool.Query(context.Background(),
		"SELECT id, number, row, is_reserved FROM seats WHERE is_reserved=$1", true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []Seats

	for rows.Next() {
		var s Seats
		err = rows.Scan(&s.Id, &s.Number, &s.Row, &s.Is_reserved)
		if err != nil {
			return nil, err
		}
		seats = append(seats, s)
	}
	return seats, nil

}
