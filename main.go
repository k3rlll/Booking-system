package main

import (
	"fmt"
	"log"
	"rest_api/db"
	"rest_api/server"
)

// type Safe struct{
// 	mtx sync.RWMutex
// }

func main() {

	// var mtx sync.RWMutex

	pool, err := db.ConnectionDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close() // закрываем только в main

	UserRepo := db.NewUserRepository(pool)
	SeatRepo := db.NewSeatRepository(pool)
	ReservationRepo := db.NewReservationRepository(pool)

	// httpHandler := &server.Http{
	// 	User:        UserRepo,
	// 	Seats:       SeatRepo,
	// 	Reservation: ReservationRepo,
	// }

	db := server.NewHttp(UserRepo, SeatRepo, ReservationRepo)
	httpServer := server.NewHTTPServer(db)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start http server:", err)
	}

}
