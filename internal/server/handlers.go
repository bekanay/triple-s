package server

import (
	"net/http"
	"strings"
)

func (s *Server) handleBuckets(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	switch r.Method {
	case http.MethodPut:
		if path == "" {
			http.Error(w, "bucket name required", http.StatusBadRequest)
			return
		}
		s.createBucket(w, r, path)

	case http.MethodGet:
		s.listBuckets(w, r)

	case http.MethodDelete:
		if path == "" {
			http.Error(w, "bucket name required", http.StatusBadRequest)
			return
		}
		s.deleteBucket(w, r, path)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) createBucket(w http.ResponseWriter, r *http.Request, name string) {
	if err := s.svc.Create(name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created bucket " + name))
}

func (s *Server) listBuckets(w http.ResponseWriter, r *http.Request) {
	names, err := s.svc.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, n := range names {
		w.Write([]byte(n + "\n"))
	}
}

func (s *Server) deleteBucket(w http.ResponseWriter, r *http.Request, name string) {
	if err := s.svc.Delete(name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("Deleted bucket " + name))
}
