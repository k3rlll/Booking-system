package db

import (
	"context"
	"rest_api/functions"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	pool *pgxpool.Pool
}

func (r *ReservationRepository) IsReserved(seatId int, pool *pgxpool.Pool) bool {

	var reserved bool

	_ = r.pool.QueryRow(context.Background(),
		"SELECT is_reserved FROM seats where id=$1", seatId).
		Scan(&reserved)

	return reserved

}

func (r *ReservationRepository) Reserve(userId int, seatId int, pool *pgxpool.Pool) error {

	if r.IsReserved(seatId, pool) {
		return functions.ErrReservationAlreadyExist
	} else {

		tx, err := r.pool.Begin(context.Background())
		if err != nil {
			return err
		}

		defer tx.Rollback(context.Background())

		_, err = tx.Exec(context.Background(),
			"INSERT INTO reservation (user_id, seat_id) VALUES ($1, $2)", userId, seatId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(context.Background(),
			"UPDATE seats SET is_reserved=true WHERE id=$1", seatId)
		if err != nil {
			return err
		}

		if err := tx.Commit(context.Background()); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReservationRepository) DeleteReservation(userId int, seatId int, pool *pgxpool.Pool) error {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	if !r.IsReserved(seatId, pool) {
		return functions.ErrReservationNotFound
	} else {

		if _, err := tx.Exec(context.Background(),
			"DELETE FROM reservation WHERE (user_id=$1 && seat_id=$2)", userId, seatId); err != nil {
			return err
		}

		if _, err = tx.Exec(context.Background(),
			"UPDATE seats SET is_reserved=false where id=$1", seatId); err != nil {
			return err
		}

	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}
	return nil

}
