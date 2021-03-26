package notifications

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
		req, err := http.NewRequest(http.MethodOptions, body.URL, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Something is wrong with the url field", http.StatusBadRequest)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Something is wrong with the url field", http.StatusBadRequest)
			return
		}
		defer res.Body.Close()
		if strings.Contains(res.Header.Get("Allow"), http.MethodPost) {
			http.Error(rw, "The supplied url exits, but does not support eater POST requests or OPTIONS requests", http.StatusBadRequest)
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
