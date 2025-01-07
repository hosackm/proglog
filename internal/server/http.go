package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewHttpServer(addr string) *http.Server {
	srv := newHttpServer()
	r := chi.NewMux()

	// consume = read
	r.Get("/consume", srv.HandleConsume)

	// produce = append
	r.Post("/produce", srv.HandleProduce)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	log *Log
}

// User wants to retrieve a record at a specified offset
func (s *httpServer) HandleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := s.log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ProduceResponse{Offset: offset}
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) HandleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rec, err := s.log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ConsumeResponse{Record: rec}
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newHttpServer() *httpServer {
	return &httpServer{
		log: NewLog(),
	}
}

// User wants to append a record to the log
type ProduceRequest struct {
	Record Record `json:"record"`
}

// The offset of the appended record
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

// User wants to retrieve a message at an offset
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

// The record at requested offset
type ConsumeResponse struct {
	Record Record `json:"record"`
}
