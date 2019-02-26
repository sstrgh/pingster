package site

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sstrgh/pingster/api/site"
)

func TestServeHTTP(t *testing.T) {

	api := &site.API{}

	t.Run("Testing empty responses for get requests to 'api/sites'", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		api.ServeHTTP(response, request)

		statusCode := response.Code
		statusCodeWant := 200

		if statusCode != statusCodeWant {
			t.Errorf("StatusCode for requests to 'api/sites' needs to be %d, got '%d',", statusCodeWant, statusCode)
		}
	})

	t.Run("Testing empty responses for post requests to 'api/sites'", func(t *testing.T) {
		emptyString := []byte(`{}`)

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(emptyString))
		response := httptest.NewRecorder()

		api.ServeHTTP(response, request)

		statusCode := response.Code
		statusCodeWant := 400

		if statusCode != statusCodeWant {
			t.Errorf("StatusCode for requests to 'api/sites' needs to be %d, got '%d',", statusCodeWant, statusCode)
		}

		validationErrors := response.Body.String()
		validationWant := "{\"errors\":[\"Invalid URL\",\"Name Required\"]}"
		if validationErrors != validationWant {
			t.Errorf("StatusCode for requests to 'api/sites' needs to be %s, got '%s',", validationWant, validationErrors)
		}
	})

	t.Run("Testing responses for proper post requests to 'api/sites'", func(t *testing.T) {
		siteReq := &site.Site{
			Name:     "google",
			Endpoint: "http://www.google.com",
		}
		requestBody, _ := json.Marshal(siteReq)

		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()

		api.ServeHTTP(response, request)

		statusCode := response.Code
		statusCodeWant := 200

		if statusCode != statusCodeWant {
			t.Errorf("StatusCode for requests to 'api/sites' needs to be %d, got '%d',",
				statusCodeWant,
				statusCode,
			)
		}

		siteResp := &site.Site{}
		json.NewDecoder(response.Body).Decode(&siteResp)

		if siteResp.Endpoint != siteReq.Endpoint {
			t.Errorf("Site endpoint needs to be %+v, got %+v,",
				siteReq.Endpoint,
				siteResp.Endpoint,
			)
		}

		// Testing to see if get request displays the newly inserted Sites
		getRequest, _ := http.NewRequest(http.MethodGet, "/", nil)
		getResponse := httptest.NewRecorder()

		api.ServeHTTP(getResponse, getRequest)

		var getRespSite map[string]site.Site
		json.NewDecoder(getResponse.Body).Decode(&getRespSite)

		_, found := getRespSite[siteReq.Endpoint]

		if !found {
			t.Errorf("Expected %+v to be here",
				getRespSite[siteReq.Endpoint],
			)
		}

	})

	t.Run("Testing responses for delete requests to 'api/sites'", func(t *testing.T) {
		// Adding google to db
		siteReq := &site.Site{
			Name:     "google",
			Endpoint: "http://www.google.com",
		}
		requestBody, _ := json.Marshal(siteReq)
		request, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		response := httptest.NewRecorder()
		api.ServeHTTP(response, request)

		// Adding facebook to db
		siteReq = &site.Site{
			Name:     "facebook",
			Endpoint: "http://www.facebook.com",
		}
		requestBody, _ = json.Marshal(siteReq)
		request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		response = httptest.NewRecorder()
		api.ServeHTTP(response, request)

		// deleting google
		deleteReq := struct {
			Endpoint string `json:"endpoint"`
		}{
			Endpoint: "http://www.google.com",
		}
		requestBody, _ = json.Marshal(deleteReq)
		request, _ = http.NewRequest(http.MethodDelete, "/", bytes.NewBuffer(requestBody))
		response = httptest.NewRecorder()
		api.ServeHTTP(response, request)

		// Checking to see if status code is 202
		statusCode := response.Code
		statusCodeWant := 202
		if statusCode != statusCodeWant {
			t.Errorf("StatusCode for requests to 'api/sites' needs to be %d, got '%d',",
				statusCodeWant,
				statusCode,
			)
		}

		// Testing to see if get request displays only the sites that are left
		getRequest, _ := http.NewRequest(http.MethodGet, "/", nil)
		getResponse := httptest.NewRecorder()

		api.ServeHTTP(getResponse, getRequest)

		var getRespSite map[string]site.Site
		json.NewDecoder(getResponse.Body).Decode(&getRespSite)

		_, found := getRespSite[deleteReq.Endpoint]

		if found {
			fmt.Printf("\n%+v\n", getRespSite)
			t.Errorf("Site %+v should not be here",
				getRespSite[deleteReq.Endpoint],
			)
		}
	})
}
