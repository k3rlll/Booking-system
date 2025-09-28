package db

import (
	"context"
	"rest_api/functions"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reserve struct {
	UserId User
	SeatId Seats
}

type ReservationRepository struct {
	pool *pgxpool.Pool
}

func NewReservationRepository(pool *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{
		pool: pool,
	}
}

func (r *ReservationRepository) IsReserved(ctx context.Context, seatId int) bool {
	var reserved = false

	if seatId == 0 {
		return reserved
	}

	_ = r.pool.QueryRow(ctx,
		"SELECT is_reserved FROM seats where id=$1", seatId).
		Scan(&reserved)

	return reserved

}

func (r *ReservationRepository) Reserve(ctx context.Context, userId int, seatId int) error {

	if r.IsReserved(ctx, seatId) {
		return functions.ErrReservationAlreadyExist
	} else {

		tx, err := r.pool.Begin(ctx)
		if err != nil {
			return err
		}

		defer tx.Rollback(ctx)

		_, err = tx.Exec(ctx,
			"INSERT INTO reservation (user_id, seat_id) VALUES ($1, $2)", userId, seatId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			"UPDATE seats SET is_reserved=true WHERE id=$1", seatId)
		if err != nil {
			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *ReservationRepository) DeleteReservation(ctx context.Context, userId int, seatId int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if !r.IsReserved(ctx, seatId) {
		return functions.ErrReservationNotFound
	} else {

		if _, err := tx.Exec(ctx,
			"DELETE FROM reservation WHERE (user_id=$1 && seat_id=$2)", userId, seatId); err != nil {
			return err
		}

		if _, err = tx.Exec(ctx,
			"UPDATE seats SET is_reserved=false where id=$1", seatId); err != nil {
			return err
		}

	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil

}
