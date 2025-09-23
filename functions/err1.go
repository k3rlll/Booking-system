package functions

import "errors"

var ErrUserAlreadyCreated = errors.New("user already created")

var ErrReservationAlreadyExist = errors.New("reservation already exist")

var ErrReservationNotFound = errors.New("reservation not found")
