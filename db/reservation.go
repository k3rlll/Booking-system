package db

import (
	"context"
	"rest_api/functions"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Reservation struct {
	Reservation_id int `json:"reservation_id"`
	UserId         int `json:"user_id"`
	SeatRow        int `json:"seat_row"`
	SeatNumber     int `json:"seat_number"`
}

type ReservationRepository struct {
	pool *pgxpool.Pool
}

func NewReservationRepository(pool *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{
		pool: pool,
	}
}

func (r *ReservationRepository) IsReserved(ctx context.Context, seatRow int, seatNumber int) bool {
	var reserved = false

	if seatNumber == 0 || seatRow == 0 {
		return reserved
	}

	_ = r.pool.QueryRow(ctx,
		"SELECT is_reserved FROM seats where row=$1 and number=$2", seatRow, seatNumber).
		Scan(&reserved)

	return reserved

}

func (r *ReservationRepository) Reserve(ctx context.Context, userId int, seatRow int, seatNumber int) (Reservation, error) {

	var res Reservation

	if r.IsReserved(ctx, seatRow, seatNumber) {
		return Reservation{}, functions.ErrReservationAlreadyExist
	} else {

		tx, err := r.pool.Begin(ctx)
		if err != nil {
			return Reservation{}, err
		}
		defer tx.Rollback(ctx)

		err = tx.QueryRow(ctx,
			"INSERT INTO reservation (user_id, seat_row, seat_number) VALUES ($1, $2, $3) RETURNING reservation_id, user_id, seat_row, seat_number",
			userId, seatRow, seatNumber).Scan(&res.Reservation_id, &res.UserId, &res.SeatRow, &res.SeatNumber)
		if err != nil {
			return Reservation{}, err
		}

		_, err = tx.Exec(ctx,
			"UPDATE seats SET is_reserved=true WHERE row=$1 AND number=$2", seatRow, seatNumber)
		if err != nil {
			return Reservation{}, err
		}

		if err := tx.Commit(ctx); err != nil {
			return Reservation{}, err
		}

		return res, nil

	}

}

func (r *ReservationRepository) DeleteReservation(ctx context.Context, user_id int, seatRow int, seatNumber int) error {
	var check_id int

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)
	_ = tx.QueryRow(ctx,
		"SELECT user_id FROM reservation WHERE (seat_row=$1 && seat_number=$2)", seatRow, seatNumber).Scan(check_id)

	if check_id != user_id {
		return functions.ErrNoPermission
	}

	if !r.IsReserved(ctx, seatRow, seatNumber) {
		return functions.ErrReservationNotFound
	} else {

		if _, err := tx.Exec(ctx,
			"DELETE FROM reservation WHERE (seat_row=$1 && seat_number=$2)", seatRow, seatNumber); err != nil {
			return err
		}

		if _, err = tx.Exec(ctx,
			"UPDATE seats SET is_reserved=false where seat_row=$1 and seat_number=$2", seatRow, seatNumber); err != nil {
			return err
		}

	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil

}
