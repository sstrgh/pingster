package site

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jasonlvhit/gocron"
)

// API handles api requests to '/sites/
type API struct{}

// Site is a representation for a site that needs to be pinged
type Site struct {
	Endpoint  string    `json:"endpoint,omitempty"`
	Name      string    `json:"name"`
	LastPing  time.Time `json:"lastPing"`
	scheduler *gocron.Scheduler
}

var db = map[string]*Site{}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		doGet(w, r)
	case http.MethodPost:
		doPost(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported method '%v' to %v\n", r.Method, r.URL)
		log.Printf("Unsupported method '%v' to %v\n", r.Method, r.URL)
	}
}
