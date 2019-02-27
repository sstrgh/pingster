package site

import (
	"encoding/json"
	"errors"
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

func validateErrors(newSite Site) struct {
	Errors []string `json:"errors"`
} {
	validationErrors := struct {
		Errors []string `json:"errors"`
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
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var newSite Site
	err := json.NewDecoder(r.Body).Decode(&newSite)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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

	pingEndpoint(&newSite)

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

func pingEndpoint(site *Site) {
	endpoint, _ := url.Parse(site.Endpoint)
	timeout, _ := time.ParseDuration("800ms")
	pingFunc := icmpPing

	err := pingFunc(endpoint, timeout)

	if err != nil {
		return
	}

	site.LastPing = time.Now()

	fmt.Printf(
		"Successfully pinged %s at %s",
		site.Endpoint,
		site.LastPing.Format("2006-01-02 15:04:05"),
	)

}

func icmpPing(endpoint *url.URL, timeout time.Duration) error {
	pinger, err := ping.NewPinger(endpoint.Hostname())

	if err != nil {
		fmt.Printf(
			"Failed to resolve host for %s at %s",
			endpoint,
			time.Now().Format("2006-01-02 15:04:05"),
		)
		return err
	}

	pinger.Count = 1
	pinger.Timeout = timeout

	// Because pinger.Run() is thread blocking we are creating channels and running pinger in a goroutine
	statsChnl := make(chan *ping.Statistics)
	pinger.OnFinish = func(s *ping.Statistics) {
		statsChnl <- s
		return
	}

	go pinger.Run()

	stats := <-statsChnl

	if stats.PacketsRecv != 1 {
		fmt.Printf(
			"Pinged %s at %s and failed to receive packets",
			endpoint,
			time.Now().Format("2006-01-02 15:04:05"),
		)
		err := errors.New("Failed to receive a response")

		return err
	}

	return nil
}

func registerPingCron(site *Site) {
	site.scheduler = gocron.NewScheduler()
	site.scheduler.Every(10).Seconds().Do(pingEndpoint, site) // Scheduled a job to run every 5mins(300 seconds)
	<-site.scheduler.Start()
}
