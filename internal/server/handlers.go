package server

import (
	"encoding/xml"
	"net/http"
	"strings"
	"triple-s/internal/storage"
)

type PutObjectResult struct {
	XMLName xml.Name `xml:"PutObjectResult"`
	ETag    string   `xml:"ETag"`
}

type PutBucketResult struct {
	XMLName xml.Name `xml:"PutBucketResult"`
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
		parts := strings.Split(path, "/")
		switch slashCounter {
		case 0:
			s.createBucket(w, r, path)
		case 1:
			s.uploadObject(w, r, parts[0], parts[1])
		default:
			resp := ErrorResponse{
				Code:    "400",
				Message: "can not use more than 2 segments in path",
			}
			writeXML(w, http.StatusBadRequest, resp)
		}

	case http.MethodGet:
		slashCounter := strings.Count(path, "/")
		parts := strings.Split(path, "/")
		switch slashCounter {
		case 0:
			s.listBuckets(w, r)
		case 1:
			s.getObject(w, r, parts[0], parts[1])
		default:
			resp := ErrorResponse{
				Code:    "400",
				Message: "can not use more than 2 segments in path",
			}
			writeXML(w, http.StatusBadRequest, resp)
		}

	case http.MethodDelete:
		if path == "" {
			http.Error(w, "bucket name required", http.StatusBadRequest)
			return
		}
		slashCounter := strings.Count(path, "/")
		parts := strings.Split(path, "/")
		switch slashCounter {
		case 0:
			s.deleteBucket(w, r, path)
		case 1:
			s.deleteObject(w, r, parts[0], parts[1])
		default:
			resp := ErrorResponse{
				Code:    "400",
				Message: "can not use more than 2 segments in path",
			}
			writeXML(w, http.StatusBadRequest, resp)
		}

	default:
		resp := ErrorResponse{
			Code:    "400",
			Message: "Method not allowed",
		}
		writeXML(w, http.StatusMethodNotAllowed, resp)
	}
}

func (s *Server) createBucket(w http.ResponseWriter, r *http.Request, name string) {
	if err := s.svc.CreateBucket(name); err != nil {
		resp := ErrorResponse{
			Code:    "409",
			Message: err.Error(),
		}
		writeXML(w, http.StatusConflict, resp)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	resp := PutObjectResult{ETag: "Created bucket " + name}
	xml.NewEncoder(w).Encode(resp)
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
	resp := PutObjectResult{ETag: "file uploaded successfully"}
	xml.NewEncoder(w).Encode(resp)
}

func (s *Server) getObject(w http.ResponseWriter, r *http.Request, bucket, object string) {
	var bytes []byte
	var contentType string
	var err error
	bytes, contentType, err = s.svc.GetObject(bucket, object)
	if err != nil {
		switch err.Error() {
		case "no bucket found", "no object found":
			resp := ErrorResponse{
				Code:    "404",
				Message: err.Error(),
			}
			writeXML(w, http.StatusNotFound, resp)
		default:
			resp := ErrorResponse{
				Code:    "500",
				Message: err.Error(),
			}
			writeXML(w, http.StatusInternalServerError, resp)
		}
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		resp := ErrorResponse{
			Code:    "500",
			Message: "failed to send file: " + err.Error(),
		}
		writeXML(w, http.StatusInternalServerError, resp)
	}
}

func (s *Server) deleteObject(w http.ResponseWriter, r *http.Request, bucket, object string) {
	if err := s.svc.DeleteObject(bucket, object); err != nil {
		resp := ErrorResponse{
			Code:    "404",
			Message: err.Error(),
		}
		writeXML(w, http.StatusNotFound, resp)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
