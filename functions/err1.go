package functions

import "errors"

var ErrUserAlreadyCreated = errors.New("user already created")

var ErrReservationAlreadyExist = errors.New("reservation already exist")

var ErrReservationNotFound = errors.New("reservation not found")

var ErrBadRequest = errors.New("request is not valid")

var ErrNoPermission = errors.New("no permission to delete this reservation")
