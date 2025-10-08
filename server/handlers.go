package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rest_api/db"
	"rest_api/functions"
	"time"
)

// HTTPHandlers
type Http struct {
	User        *db.UserRepository
	Seats       *db.SeatsRepository
	Reservation *db.ReservationRepository
}

func NewHttp(user *db.UserRepository, seats *db.SeatsRepository, reservation *db.ReservationRepository) *Http {
	return &Http{
		User:        user,
		Seats:       seats,
		Reservation: reservation,
	}
}

/*
pattern: /users
method:  POST
info:    JSON in HTTP request body

succeed:
  - status code:   201 OK
  - response body: JSON represent request to creare new user

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *Http) HandlerNewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u db.User
	var u_res db.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now()}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	u_res, err := h.User.NewUser(ctx, u.Name, u.Email)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		switch {
		case errors.Is(err, functions.ErrUserAlreadyCreated):
			http.Error(w, errDTO.ToString(), http.StatusConflict)
			return
		case errors.Is(err, functions.ErrBadRequest):
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		default:
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	// ID, err := h.User.GetUserID(ctx, u.Email)
	// if err != nil {
	// 	errDTO := functions.ErrDTO{
	// 		Error: err,
	// 		Time:  time.Now(),
	// 	}
	// 	http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	// }

	b, err := json.MarshalIndent(u_res, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
	}

}

/*
pattern: /users
method:  GET
info:    JSON in HTTP request body

succeed:
  - status code:   200 OK
  - response body: JSON represent request to get list

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *Http) HandlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u db.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	res, err := h.User.GetUserByID(ctx, u.Id)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		if err == functions.ErrBadRequest {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
			return
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			return
		}
	}

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
	}
}

/*
pattern: /seats
method:  GET
info:    JSON in HTTP request body

succeed:
  - status code:   200 OK
  - response body: JSON represent request to get list

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *Http) HandlerGetAllSeats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.Seats.GetAllSeats(ctx)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
	}

}

/*
pattern: /seats?free==true
method:  GET
info:    JSON in HTTP request body

succeed:
  - status code:   200 OK
  - response body: JSON represent request to get list

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *Http) HandlerGetFreeSeats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.Seats.GetFreeSeats(ctx)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	// if err := json.NewEncoder(w).Encode(res); err != nil {
	// 	errDTO := functions.ErrDTO{
	// 		Error: err.Error(),
	// 		Time:  time.Now(),
	// 	}
	// 	http.Error(w, errDTO.ToString(), http.StatusBadRequest)

	// }

	b, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
	}

	// w.WriteHeader(http.StatusOK)
}

/*
pattern: /seats?free==false
method:  GET
info:    JSON in HTTP request body

succeed:
  - status code:   200 OK
  - response body: JSON represent request to get list

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *Http) HandlerGetReservedSeats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.Seats.GetReservedSeats(ctx)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)

	}

	w.WriteHeader(http.StatusOK)
}

/*
pattern: /reservation
method:  GET
info:    JSON in HTTP request body

succeed:
  - status code:   200 OK
  - response body: JSON represent request to get list

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *Http) HandlerIsReserved(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var s db.Seats
	var res bool
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	if s.Row == 0 || s.Number == 0 {
		http.Error(w, functions.ErrBadRequest.Error(), http.StatusBadRequest)
	}

	res = h.Reservation.IsReserved(ctx, s.Row, s.Number)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

/*
pattern: /reservation
method:  PATCH
info:    JSON in HTTP request body

succeed:
  - status code:   202 OK
  - response body: JSON represent request to get list

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *Http) HandlerReserve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var res db.Reservation

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
	}

	result, err := h.Reservation.Reserve(ctx, res.UserId, res.SeatRow, res.SeatNumber)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		if errors.Is(err, functions.ErrReservationAlreadyExist) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		}
	}
	b, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
	}
}

/*
pattern: /reservation
method:  DELETE
info:    ----------------

succeed:
  - status code:   202 OK
  - response body: -------------

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *Http) HandlerDeleteReservation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var res db.Reservation

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
	}

	if err := h.Reservation.DeleteReservation(ctx, res.UserId, res.SeatRow, res.SeatNumber); err != nil {
		errDTO := functions.ErrDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		if errors.Is(err, functions.ErrReservationNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		}
	}
}
