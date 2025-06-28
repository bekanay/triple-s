package bucket

import (
	"errors"
	"os"
	"path/filepath"
	"time"
	"triple-s/internal/validate"
)

type Service interface {
	Create(name string) error
	List() ([]string, error)
	Delete(name string) error
}

type service struct {
	baseDir string
}

func NewService(baseDir string) (Service, error) {
	// ensure CSV file exists, etc.
	if err := initBucketsCSV(baseDir); err != nil {
		return nil, err
	}
	return &service{baseDir: baseDir}, nil
}

func (s *service) Create(name string) error {
	if err := validate.Name(name); err != nil {
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

func (s *service) List() ([]string, error) {
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

func (s *service) Delete(name string) error {
	// remove dir, remove CSV row, etc.
	return os.RemoveAll(filepath.Join(s.baseDir, name))
}
