package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Theater struct {
	name string
}

type Movie struct {
	name string
}

type Ticket struct {
	seansId int
	dateStart time.Time
	dateEnd time.Time
	price int
	movieName string
	movieTheater string
	availableSeats string
}

type DB interface {
	GetTickets(getTicketsRequest *getTicketsRequest) ([]Ticket, error)
	GetTheaterMovies(getTheaterMoviesRequest *getTheaterMoviesRequest) ([]Movie, error)
	GetMovieTheaters(getTheatersRequest *getMovieTheatersRequest) ([]Theater, error)
}

type HttpServer struct {
	httpServer *http.Server
	db DB
}

func NewHttpServer(addr string, db DB) *HttpServer {
	server := HttpServer{ db: db }

	mux := http.NewServeMux()
	mux.HandleFunc("/tickets", server.GetTickets)
	mux.HandleFunc("/theaters", server.GetMovieTheaters)
	mux.HandleFunc("/movies", server.GetTheaterMovies)

	httpServer := http.Server{Addr: addr, Handler: mux}
	server.httpServer = &httpServer

	return &server
}

func (c *HttpServer) ListenAndServe() error {
	return c.httpServer.ListenAndServe()
}

type getTicketsRequest struct {
	dateFrom time.Time
	dateTo time.Time
	priceFrom int
	priceTo int
	movieTheater string
	name string
}

type getTheaterMoviesRequest struct {
	dateFrom time.Time
	dateTo time.Time
}

type getMovieTheatersRequest struct {
	name string
}

func (s *HttpServer) GetTickets(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(400)
		return
	}

	query := r.URL.Query()

	queryDateFrom := query.Get("dateFrom")
	if queryDateFrom == "" {
		rw.WriteHeader(400)
		return
	}
	dateFrom, err := time.Parse(time.RFC822, queryDateFrom)
	if err != nil {
		rw.WriteHeader(400)
		return
	}
	queryDateTo := query.Get("dateTo")
	if queryDateTo == "" {
		rw.WriteHeader(400)
		return
	}
	dateTo, err := time.Parse(time.RFC822, queryDateTo)
	if err != nil {
		rw.WriteHeader(400)
		return
	}


	queryPriceFrom := query.Get("priceFrom")
	if queryPriceFrom == "" {
		rw.WriteHeader(400)
		return
	}
	priceFrom, err := strconv.Atoi(queryPriceFrom)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	queryPriceTo := query.Get("priceTo")
	if queryPriceTo == "" {
		rw.WriteHeader(400)
		return
	}
	priceTo, err := strconv.Atoi(queryPriceTo)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	queryMovieTheater := query.Get("movieTheater")
	if queryMovieTheater == "" {
		rw.WriteHeader(400)
		return
	}

	queryName := query.Get("name")
	if queryName == "" {
		rw.WriteHeader(400)
		return
	}

	getTicketsRequest := getTicketsRequest{
		dateFrom:     dateFrom,
		dateTo:       dateTo,
		priceFrom:    priceFrom,
		priceTo:      priceTo,
		movieTheater: queryMovieTheater,
		name:         queryName,
	}

	theaters, err := s.db.GetTickets(&getTicketsRequest)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	theatersJson, err := json.Marshal(&theaters)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
	rw.Write(theatersJson)
}

func (s *HttpServer) GetTheaterMovies(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(400)
		return
	}

	query := r.URL.Query()

	queryDateFrom := query.Get("dateFrom")
	if queryDateFrom == "" {
		rw.WriteHeader(400)
		return
	}
	dateFrom, err := time.Parse(time.RFC822, queryDateFrom)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	queryDateTo := query.Get("dateTo")
	if queryDateTo == "" {
		rw.WriteHeader(400)
		return
	}
	dateTo, err := time.Parse(time.RFC822, queryDateTo)
	if err != nil {
		rw.WriteHeader(400)
		return
	}

	getTheaterMoviesRequest := getTheaterMoviesRequest{
		dateFrom:     dateFrom,
		dateTo:       dateTo,
	}

	theaters, err := s.db.GetTheaterMovies(&getTheaterMoviesRequest)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	theatersJson, err := json.Marshal(&theaters)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
	rw.Write(theatersJson)
}

func (s *HttpServer) GetMovieTheaters(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rw.WriteHeader(400)
		return
	}

	query := r.URL.Query()

	name := query.Get("name")
	if name == "" {
		rw.WriteHeader(400)
		return
	}

	getMovieTheatersRequest := getMovieTheatersRequest{
		name: name,
	}

	theaters, err := s.db.GetMovieTheaters(&getMovieTheatersRequest)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	theatersJson, err := json.Marshal(&theaters)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
	rw.Write(theatersJson)
}