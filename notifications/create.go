package notifications

import (
	"assignment-2/corona"
	"encoding/json"
	"net/http"
)

// requestBody is the body of the incoming request to be parsed and registered.
type requestBody struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}

// responseBody is the body of the response sent back from the webhook creation endpoint.
type responseBody struct {
	ID string `json:"id"`
}

func NewCreateHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Validation:
		// - Send OPTIONS request to provided url and check if is exists and accepts POST requests
		// - Check the field is one of the enumerated options
		// - Check the trigger is one of the enumerated options

		// Decode the request body into a struct
		var body requestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(rw, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Send an OPTIONS request to the supplied url in order to check:
		// 1. That the url exits.
		// 2. That the url accepts POST requests.
		status := corona.GetStatusOf(body.URL)
		if !corona.StatusIs2XX(status) {
			http.Error(rw, "There is something wrong with the url field", http.StatusBadRequest)
			return
		}

		// Check if the field is valid
		if body.Field != FieldStringency && body.Field != FieldConfirmed {
			http.Error(rw, "The field supplied does not exits", http.StatusBadRequest)
			return
		}

		// Check if the trigger is valid
		if body.Trigger != TriggerOnChange && body.Trigger != TriggerOnTimeout {
			http.Error(rw, "The trigger supplied does not exits", http.StatusBadRequest)
			return
		}

		// Now actually create / register the webhook

		http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
