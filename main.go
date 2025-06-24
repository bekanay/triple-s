package main

import (
	"encoding/csv"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Incorrect usage. Type triple-s --help")
	}
	tripleCmd := flag.NewFlagSet("triplet", flag.ExitOnError)
	portNumber := tripleCmd.String("port", "", "port number on which server is running")
	isHelp := tripleCmd.Bool("help", false, "usage of triplet-s program")
	directory := tripleCmd.String("dir", "", "directory of storage")
	tripleCmd.Parse(os.Args[1:])

	if *isHelp {
		log.Fatal("usage")
	}

	err := initBucketsCSV(*directory)
	if err != nil {
		log.Fatalf("Failed to initialize buckets.csv: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", route)
	log.Println("Starting server on :", *portNumber)
	err = http.ListenAndServe(*portNumber, mux)
	log.Fatal(err)
}

func route(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if r.Method == http.MethodPut {
		if len(path) > 1 {
			path = path[1:]
			bucketName := path
			handleCreateBucket(w, r, bucketName)
			return
		}
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
	}

	if r.Method == http.MethodDelete {
	}
}

func handleCreateBucket(w http.ResponseWriter, r *http.Request, bucketName string) {
	if len(bucketName) < 3 || len(bucketName) > 63 {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	for i := 0; i < len(bucketName); i++ {
		if bucketName[i] >= 'a' && bucketName[i] <= 'z' {
			continue
		}
		if bucketName[i] >= '0' && bucketName[i] <= '9' {
			continue
		}
		if bucketName[i] == '-' || bucketName[i] == '.' {
			continue
		}
		http.Error(w, http.StatusText(400), 400)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func handleListBuckets(w http.ResponseWriter, r *http.Request) {
}

func handleDeleteBucket(w http.ResponseWriter, r *http.Request) {
}

func initBucketsCSV(dir string) error {
	filePath := filepath.Join(dir, "buckets.csv")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		headers := []string{"Name", "CreationTime", "LastModifiedTime"}
		if err := writer.Write(headers); err != nil {
			return err
		}
	}
	return nil
}
