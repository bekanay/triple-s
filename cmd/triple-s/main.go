package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"triple-s/internal/server"
)

const helpText = `Simple Storage Service.

Usage:
    triple-s [-port <N>] [-dir <S>]
    triple-s --help

Options:
  --help       Show this screen.
  --port N     Port number
  --dir S      Path to the directory
`

func main() {
	port := flag.String("port", "8080", "port to listen on")
	dir := flag.String("dir", "./data", "storage directory")
	help := flag.Bool("help", false, "usage of triplet-s program")
	flag.Parse()

	if *help {
		fmt.Fprint(os.Stderr, helpText)
		return
	}

	srv, err := server.New(*dir)
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}

	log.Printf("listening on :%s …", *port)
	if err := srv.Run(":" + *port); err != nil {
		log.Fatal(err)
	}
}
