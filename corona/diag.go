package corona

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
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

// countWebhooks returns the number of registered webhooks.
func countWebhooks(fs *firestore.Client) (int, error) {
	// Get refs to all the documents in the collection
	// FIXME: Use the WebhookCollection constant instead, but it's a cyclic dependency and I'm at T-1, so f it
	docrefs, err := fs.Collection("webhooks").DocumentRefs(context.Background()).GetAll()
	if err != nil {
		return 0, err
	}

	return len(docrefs), nil
}

// NewDiagHandler returns a handler function for the diagnostic endpoint.
func NewDiagHandler(fs *firestore.Client, startTime time.Time) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		uptime := int(time.Since(startTime).Seconds())

		registered, err := countWebhooks(fs)
		if err != nil {
			log.Println("Error while counting webhooks:", err.Error())
			registered = -1 // Just return a nonsense count, which more useful in this particular case
		}

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
