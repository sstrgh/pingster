package site

import (
	"encoding/json"
	"net/http"
)

func doDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	var siteToDelete struct {
		Endpoint string `json:"endpoint"`
	}
	err := json.NewDecoder(r.Body).Decode(&siteToDelete)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, found := db[siteToDelete.Endpoint]

	if found {
		val.scheduler.Clear()
	}

	delete(db, siteToDelete.Endpoint)

	w.WriteHeader(http.StatusAccepted)
	je := json.NewEncoder(w)
	je.Encode(db)
}
