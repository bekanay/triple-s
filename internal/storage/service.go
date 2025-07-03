package storage

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Storage interface {
	CreateBucket(name string) error
	ListBuckets() ([]string, error)
	DeleteBucket(name string) error

	UploadObject()
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

func (s *service) UploadObject() {
}

func (s *service) GetObject() {
}

func (s *service) DeleteObject() {
}
