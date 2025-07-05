package storage

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Storage interface {
	CreateBucket(name string) error
	ListBuckets() ([]string, error)
	DeleteBucket(name string) error

	UploadObject(bucket string, r io.Reader, object, contentType string) error
	GetObject()
	DeleteObject()
}

type service struct {
	baseDir string
}

func NewService(baseDir string) (Storage, error) {
	// ensure CSV file exists, etc.
	if err := initBucketsCSV(baseDir); err != nil {
		return nil, err
	}
	return &service{baseDir: baseDir}, nil
}

func (s *service) CreateBucket(name string) error {
	if err := Name(name); err != nil {
		return err
	}
	path := filepath.Join(s.baseDir, name)
	if _, err := os.Stat(path); err == nil {
		return errors.New("bucket already exists")
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return appendBucketToCSV(s.baseDir, name, time.Now())
}

func (s *service) ListBuckets() ([]string, error) {
	// e.g. read directories under baseDir
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func (s *service) DeleteBucket(name string) error {
	if err := Name(name); err != nil {
		return err
	}
	path := filepath.Join(s.baseDir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New(("bucket not found"))
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return errors.New("bucket not empty")
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return removeBucketFromCSV(s.baseDir, name)
}

func (s *service) UploadObject(bucket string, r io.Reader, object, contentType string) error {
	path := filepath.Join(s.baseDir, bucket)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("no bucket found")
	}
	if strings.TrimSpace(object) == "" {
		return errors.New("invalid object key")
	}
	metaFile := filepath.Join(path, "objects.csv")
	if err := ensureObjectsCSV(metaFile); err != nil {
		return err
	}
	path = filepath.Join(path, object)
	f, err := os.Create(path)
	if err != nil {
		return errors.New("error while creating object file")
	}
	defer f.Close()

	bytesWritten, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	return upsertObjectMetadata(metaFile, object, bytesWritten, contentType, now)
}

func upsertObjectMetadata(path, key string, size int64, ctype, modTime string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	reader := csv.NewReader(f)
	recs, err := reader.ReadAll()
	f.Close()
	if err != nil {
		return err
	}
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.Write(recs[0])
	updated := false
	for _, row := range recs[1:] {
		if row[0] == key {
			w.Write([]string{key, strconv.FormatInt(size, 10), ctype, modTime})
			updated = true
		} else {
			w.Write(row)
		}
	}
	if !updated {
		w.Write([]string{key, strconv.FormatInt(size, 10), ctype, modTime})
	}
	w.Flush()
	return w.Error()
}

func ensureObjectsCSV(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		return w.Write([]string{"ObjectKey", "Size", "ContentType", "LastModified"})
	}
	return nil
}

func (s *service) GetObject() {
}

func (s *service) DeleteObject() {
}
