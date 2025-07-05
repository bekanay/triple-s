package server

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"triple-s/internal/storage"
)

type PutObjectResult struct {
	XMLName xml.Name `xml:"PutObjectResult"`
	ETag    string   `xml:"ETag"`
}

func (s *Server) handleBuckets(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	switch r.Method {
	case http.MethodPut:
		if path == "" {
			http.Error(w, "bucket name required", http.StatusBadRequest)
			return
		}
		slashCounter := strings.Count(path, "/")
		fmt.Println(slashCounter)
		parts := strings.Split(path, "/")
		switch slashCounter {
		case 0:
			s.createBucket(w, r, path)
		case 1:
			s.uploadObject(w, r, parts[0], parts[1])
		default:
			http.Error(w, "can not use more than 2 segments in path", http.StatusBadRequest)
		}

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
	if err := s.svc.CreateBucket(name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created bucket " + name))
}

func (s *Server) listBuckets(w http.ResponseWriter, r *http.Request) {
	metas, err := storage.ReadAllMetadata(s.baseDir)
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
	if err := s.svc.DeleteBucket(name); err != nil {
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

func (s *Server) uploadObject(w http.ResponseWriter, r *http.Request, bucket, object string) {
	contentType := r.Header.Get("Content-Type")

	err := s.svc.UploadObject(bucket, r.Body, object, contentType)
	if err != nil {
		switch err.Error() {
		case "no bucket found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "invalid object key":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	resp := PutObjectResult{ETag: ""}
	xml.NewEncoder(w).Encode(resp)
}
