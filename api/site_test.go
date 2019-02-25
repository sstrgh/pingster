package site

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
		validationWant := "{\"Errors\":[\"Invalid URL\",\"Name Required\"]}"
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

		if !cmp.Equal(siteReq, siteResp, cmpopts.IgnoreUnexported(site.Site{})) {
			t.Errorf("StatusCode for requests to 'api/sites' needs to be %+v, got %+v,",
				siteReq,
				siteResp,
			)
		}

		// Testing to see if get request displays the newly inserted Sites
		getRequest, _ := http.NewRequest(http.MethodGet, "/", nil)
		getResponse := httptest.NewRecorder()

		api.ServeHTTP(getResponse, getRequest)

		var getRespSite map[string]site.Site
		json.NewDecoder(getResponse.Body).Decode(&getRespSite)

		if !cmp.Equal(*siteReq, getRespSite[siteReq.Endpoint], cmpopts.IgnoreUnexported(site.Site{})) {
			t.Errorf("StatusCode for requests to 'api/sites' needs to be %+v, got %+v,",
				*siteReq,
				getRespSite[siteReq.Endpoint],
			)
		}

	})

}
