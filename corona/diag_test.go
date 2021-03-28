package corona

import (
	fs "assignment-2/firestore"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestDiagEndpoint tests that diag responds with the expected status code, body, and content-type.
// It also test that all the third party apis return with the expected status code.
func TestDiagEndpoint(t *testing.T) {
	// Initialize a firestore client
	fs := fs.NewFirestoreClient()
	defer fs.Close()

	req, err := http.NewRequest(http.MethodGet, "/exchange/v1/diag", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	rr := httptest.NewRecorder()
	handler := NewDiagHandler(fs, time.Now())
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK, "Status code should be 200 status ok")

	var body diag
	err = json.NewDecoder(rr.Body).Decode(&body)
	if err != nil {
		t.Fatal(err.Error())
	}

	// NOTE: CovidTrackerAPI returns "204 no content" for options requests for some reason.
	assert.True(t, StatusIs2XX(body.CovidTrackerAPI), "Status code is not 2XX")
	assert.True(t, StatusIs2XX(body.MMediaGroupAPI), "Status code is not 2XX")
	assert.True(t, StatusIs2XX(body.RestCountriesAPI), "Status code is not 2XX")
}
