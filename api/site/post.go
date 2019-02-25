package site

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/sparrc/go-ping"
)

// isValidUrl tests a string to determine if it is a url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)

	if err != nil {
		return false
	}

	return true
}

func validateErrors(newSite Site) struct{ Errors []string } {
	validationErrors := struct {
		Errors []string
	}{}

	if !isValidURL(newSite.Endpoint) {
		validationErrors.Errors = append(validationErrors.Errors, "Invalid URL")
	}

	if newSite.Name == "" {
		validationErrors.Errors = append(validationErrors.Errors, "Name Required")
	}

	return validationErrors
}

func doPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newSite Site

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&newSite)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	validationErrors := validateErrors(newSite)

	if len(validationErrors.Errors) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		b, errors := json.Marshal(validationErrors)

		if errors != nil {
			return
		}

		fmt.Fprintf(w, string(b[:]))
		return
	}

	_, found := db[newSite.Endpoint]

	if found {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{ "error": "Site already exists!" }`)
		return
	}

	pingSite(&newSite)

	db[newSite.Endpoint] = &newSite
	go registerPingCron(&newSite)

	respSite := Site{
		Endpoint: newSite.Endpoint,
		Name:     newSite.Name,
		LastPing: newSite.LastPing,
	}

	je := json.NewEncoder(w)
	je.Encode(respSite)
}

func pingSite(site *Site) {
	endpoint, _ := url.Parse(site.Endpoint)
	pinger, err := ping.NewPinger(endpoint.Hostname())

	if err != nil {
		return
	}

	pinger.Count = 1
	pinger.Timeout = 800000000 // 800 milliseconds

	// Because pinger.Run() is thread blocking we are creating channels and running pinger in a goroutine
	statsChnl := make(chan *ping.Statistics)
	pinger.OnFinish = func(s *ping.Statistics) {
		statsChnl <- s
		return
	}

	go pinger.Run()

	stats := <-statsChnl

	if stats.PacketsRecv != 1 {
		return
	}

	site.LastPing = time.Now()

	fmt.Printf(
		"Pinged %s at %s and received %d packets",
		site.Endpoint,
		site.LastPing.Format("2006-01-02 15:04:05"),
		stats.PacketsRecv,
	)
}

func registerPingCron(site *Site) {
	gocron.Every(10).Seconds().Do(pingSite, site)
	<-gocron.Start()
}
