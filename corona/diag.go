package corona

import (
	"encoding/json"
	"net/http"
	"time"
)

// diag is the response from the diagnostic interface.
type diag struct {
	MMediaGroupAPI   int    `json:"mmediagroupapi"`
	CovidTrackerAPI  int    `json:"covidtrackerapi"`
	RestCountriesAPI int    `json:"restcountriesapi"`
	Registered       int    `json:"registered"`
	Version          string `json:"version"`
	Uptime           int    `json:"uptime"`
}

// NewDiagHandler returns a handler function for the diagnostic endpoint
func NewDiagHandler(registered int, startTime time.Time) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uptime := int(time.Since(startTime).Seconds())

		response := diag{
			GetStatusOf(MMediaGroupAPIRootPath + "/cases"),
			GetStatusOf(CovidTrackerAPIRootPath[0 : len(CovidTrackerAPIRootPath)-3]),
			GetStatusOf(RestCountriesRootPath),
			registered,
			Version,
			uptime,
		}

		_ = json.NewEncoder(rw).Encode(response)
	}
}
