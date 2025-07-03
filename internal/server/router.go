package server

import (
	"net/http"
	"triple-s/internal/storage"
)

type Server struct {
	svc     storage.Storage
	mux     *http.ServeMux
	baseDir string
}

func New(dataDir string) (*Server, error) {
	svc, err := storage.NewService(dataDir)
	if err != nil {
		return nil, err
	}

	s := &Server{
		svc:     svc,
		mux:     http.NewServeMux(),
		baseDir: dataDir,
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
