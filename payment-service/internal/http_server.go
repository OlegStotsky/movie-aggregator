package internal

import (
	"encoding/json"
	"net/http"
)

type DB interface {
	AddPayment(addPaymentRequest *addPaymentRequest) error
}

type HttpServer struct {
	httpServer *http.Server
	db DB
}

func NewHttpServer(addr string, db DB) *HttpServer {
	server := HttpServer{ db: db }

	mux := http.NewServeMux()
	mux.HandleFunc("/payment", server.AddPayment)

	httpServer := http.Server{Addr: addr, Handler: mux}
	server.httpServer = &httpServer

	return &server
}

func (c *HttpServer) ListenAndServe() error {
	return c.httpServer.ListenAndServe()
}

type addPaymentRequest struct {
	clientId string
	seatId int
}

func (s *HttpServer) AddPayment(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(400)
		return
	}

	var addPaymentRequest addPaymentRequest

	err := json.NewDecoder(r.Body).Decode(&addPaymentRequest)
	if err != nil {
		rw.WriteHeader(400)
		return
	}


	err = s.db.AddPayment(&addPaymentRequest)
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
}