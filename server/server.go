package server

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	Http *Http
}

func NewHTTPServer(HTTPhandler *Http) *HTTPServer {
	return &HTTPServer{
		Http: HTTPhandler,
	}
}

func (s *HTTPServer) StartServer() error {

	router := mux.NewRouter()

	router.Path("/users").Methods("POST").HandlerFunc(s.Http.HandlerNewUser)
	router.Path("/users").Methods("GET").HandlerFunc(s.Http.HandlerGetUserByID)
	router.Path("/seats").Methods("GET").Queries("free", "true").HandlerFunc(s.Http.HandlerGetFreeSeats)
	router.Path("/seats").Methods("GET").Queries("free", "false").HandlerFunc(s.Http.HandlerGetReservedSeats)
	router.Path("/seats").Methods("GET").HandlerFunc(s.Http.HandlerGetAllSeats)
	router.Path("/reservation").Methods("GET").HandlerFunc(s.Http.HandlerIsReserved)
	router.Path("/reservation").Methods("PATCH").HandlerFunc(s.Http.HandlerReserve)
	router.Path("/reservation").Methods("DELETE").HandlerFunc(s.Http.HandlerDeleteReservation)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}
