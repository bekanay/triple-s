package server

import (
	"net/http"
	"triple-s/internal/bucket"
)

type Server struct {
	svc bucket.Service
	mux *http.ServeMux
}

func New(dataDir string) (*Server, error) {
	svc, err := bucket.NewService(dataDir)
	if err != nil {
		return nil, err
	}

	s := &Server{
		svc: svc,
		mux: http.NewServeMux(),
	}
	s.routes()
	return s, nil
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleBuckets)
}
