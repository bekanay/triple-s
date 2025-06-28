package bucket

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Metadata struct {
	Name             string
	CreationTime     string
	LastModifiedTime string
}

func initBucketsCSV(dir string) error {
	file := filepath.Join(dir, "buckets.csv")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()
		w := csv.NewWriter(f)
		defer w.Flush()
		return w.Write([]string{"Name", "CreationTime", "LastModifiedTime"})
	}
	return nil
}

func appendBucketToCSV(dir, name string, t time.Time) error {
	file := filepath.Join(dir, "buckets.csv")
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	record := []string{name, t.Format(time.RFC3339), t.Format(time.RFC3339)}
	if err := w.Write(record); err != nil {
		return fmt.Errorf("csv write: %w", err)
	}
	return nil
}

func ReadAllMetadata(dir string) ([]Metadata, error) {
	f, err := os.Open(filepath.Join(dir, "buckets.csv"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rdr := csv.NewReader(f)
	rows, err := rdr.ReadAll()
	if err != nil {
		return nil, err
	}

	var metas []Metadata
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}
		metas = append(metas, Metadata{
			Name:             row[0],
			CreationTime:     row[1],
			LastModifiedTime: row[2],
		})
	}
	return metas, nil
}

func removeBucketFromCSV(dir, name string) error {
	filePath := filepath.Join(dir, "buckets.csv")
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	rdr := csv.NewReader(f)
	records, err := rdr.ReadAll()
	defer f.Close()

	if err != nil {
		return err
	}
	out := [][]string{}

	for i, rec := range records {
		if i == 0 || rec[0] != name {
			out = append(out, rec)
		}
	}
	tmp := filepath.Join(dir, "buckets.tmp")
	f2, err := os.Create(tmp)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f2)
	if err := w.WriteAll(out); err != nil {
		f2.Close()
		return err
	}
	f2.Close()
	return os.Rename(tmp, filePath)
}
