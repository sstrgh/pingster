package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sstrgh/pingster/api/site"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	// Handles web files
	http.Handle("/", http.FileServer(http.Dir("web")))

	// register handler for '/sites'
	http.Handle("/api/sites", &site.API{})

	// Ensures that the server is up and running

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
