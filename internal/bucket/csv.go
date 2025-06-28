package bucket

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
