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

type Http struct {
	User        *db.UserRepository
	Seats       *db.SeatsRepository
	Reservation *db.ReservationRepository
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

func (h *Http) NewUser(w http.ResponseWriter, r http.Request) {
	ctx := r.Context()

	var u db.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now()}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
	}

	if err := h.User.NewUser(ctx, u.Name, u.Email); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		switch {
		case errors.Is(err, functions.ErrUserAlreadyCreated):
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		case errors.Is(err, functions.ErrBadRequest):
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		default:
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
	}

	ID, err := h.User.GetUserID(ctx, u.Email)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	b, err := json.Marshal(ID)
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

func (h *Http) GetUserByID(w http.ResponseWriter, r http.Request) {
	ctx := r.Context()

	var u db.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
	}

	res, err := h.User.GetUserByID(ctx, u.Id)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		if err == functions.ErrBadRequest {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
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

func (h *Http) GetAllSeats(w http.ResponseWriter, r http.Request) {
	ctx := r.Context()

	res, err := h.Seats.GetAllSeats(ctx)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

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

func (h *Http) GetFreeSeats(w http.ResponseWriter, r http.Request) {
	ctx := r.Context()

	res, err := h.Seats.GetFreeSeats(ctx)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)

	}

	w.WriteHeader(http.StatusOK)
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

func (h *Http) GetReservedSeats(w http.ResponseWriter, r http.Request) {
	ctx := r.Context()

	res, err := h.Seats.GetReservedSeats(ctx)
	if err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
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

func (h *Http) IsReserved(w http.ResponseWriter, r http.Request) {
	ctx := r.Context()

	var IdSeat int
	var res bool
	if err := json.NewDecoder(r.Body).Decode(&IdSeat); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
	}

	if IdSeat == 0 {
		http.Error(w, functions.ErrBadRequest.Error(), http.StatusBadRequest)
	}

	res = h.Reservation.IsReserved(ctx, IdSeat)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
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

func (h *Http) Reserve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var res db.Reserve

	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		errDTO := functions.ErrDTO{
			Error: err,
			Time:  time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
	}

	h.Reservation.Reserve(ctx)
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

func (h *Http) DeleteReservation(w http.ResponseWriter, r *http.Request) {

}
