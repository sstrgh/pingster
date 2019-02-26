package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sstrgh/pingster/api/site"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("web")))

	// register handler for '/sites'
	http.Handle("/api/sites", &site.API{})

	fmt.Println("Starting server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
