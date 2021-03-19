package corona

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

// TestDiagEndpoint tests that diag responds with the expected status code, body, and content-type.
// It also test that all the third party apis return with the expected status code.
func TestDiagEndpoint(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/exchange/v1/diag", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	rr := httptest.NewRecorder()
	handler := NewDiagHandler(0, time.Now())
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK, "Status code should be 200 status ok")

	var body diag
	err = json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Fatal(err.Error())
	}

	// NOTE: CovidTrackerAPI returns "204 no content" for options requests for some reason.
	assert.Equal(t, body.CovidTrackerAPI, http.StatusNoContent, "Status code should be 204 status no content")
	assert.Equal(t, body.MMediaGroupAPI, http.StatusOK, "Status code should be 200 status ok")
	assert.Equal(t, body.RestCountriesAPI, http.StatusOK, "Status code should be 200 status ok")
}
