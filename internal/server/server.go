package server

import (
	"net/http"
	"time"
)

func (s *Server) StartServer() {
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

type Server struct {
	*http.Server
}

func NewServer(h Handle) Server {
	mux := http.NewServeMux()
	mux.Handle("/sync", http.HandlerFunc(h.SynchroniseBookings))
	return Server{&http.Server{
		Addr:              ":8080",
		Handler:           mux,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    16 * 1024,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}}
}
