package main

import (
	"github.com/olegstotsky/movie-aggregator/theaters-service/internal"
)

func main() {
	httpServer := internal.NewHttpServer("localhost:2900", internal.DBImpl{})

	httpServer.ListenAndServe()
}
