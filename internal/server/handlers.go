package server

import (
	"encoding/xml"
	"net/http"
	"strings"
	"triple-s/internal/bucket"
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
	metas, err := bucket.ReadAllMetadata(s.baseDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := listBucketsResultXML{}
	for _, m := range metas {
		resp.Buckets = append(resp.Buckets, bucketXML{
			Name:             m.Name,
			CreationTime:     m.CreationTime,
			LastModifiedTime: m.LastModifiedTime,
		})
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) deleteBucket(w http.ResponseWriter, r *http.Request, name string) {
	if err := s.svc.Delete(name); err != nil {
		switch err.Error() {
		case "bucket not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "bucket not empty":
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
