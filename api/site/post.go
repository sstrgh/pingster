package site

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	db[newSite.Endpoint] = &newSite

	respSite := Site{
		Endpoint: newSite.Endpoint,
		Name:     newSite.Name,
		LastPing: newSite.LastPing,
	}

	je := json.NewEncoder(w)
	je.Encode(respSite)
}
