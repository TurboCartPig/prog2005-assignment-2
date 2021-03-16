package corona

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// diag is the response from the diagnostic interface.
type diag struct {
	MMediaGroupAPI  int    `json:"mmediagroupapi"`
	CovidTrackerAPI int    `json:"covidtrackerapi"`
	Registered      int    `json:"registered"`
	Version         string `json:"version"`
	Uptime          int    `json:"uptime"`
}

// getStatusOf returns the status code of a head request to the root path of a remote.
func getStatusOf(addr string) int {
	req, err := http.NewRequest(http.MethodOptions, addr, nil)
	if err != nil {
		log.Printf("Options request failed with: %s", err.Error())
		return http.StatusBadRequest // Assume I did something wrong, all other errors should be "successful"
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Options request failed with: %s", err.Error())
		return http.StatusBadRequest // Assume I did something wrong, all other errors should be "successful"
	}
	res.Body.Close()
	return res.StatusCode
}

// NewDiagHandler returns a handler function for the diagnostic endpoint
func NewDiagHandler(registered int, startTime time.Time) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uptime := int(time.Since(startTime).Seconds())

		response := diag{
			getStatusOf(MMediaGroupAPIRootPath + "/cases"),
			getStatusOf(CovidTrackerAPIRootPath[0 : len(CovidTrackerAPIRootPath)-3]),
			registered,
			Version,
			uptime,
		}

		json.NewEncoder(rw).Encode(response)
	}
}
