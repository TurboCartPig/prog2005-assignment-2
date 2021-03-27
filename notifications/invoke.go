package notifications

import (
	"assignment-2/corona"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"net/http"
	"sync"
)

// There are two types of webhooks
// 1. ON_TIMEOUT, which is easy
// 2. ON_CHANGE, but how to know if something changed?
// I will make two queues with ON_CHANGE webhooks. One for "stringency" and one for "confirmed".
// In each queue I wait for the first timeout.
// The I respond to all the webhooks if the data has changed from last time a WEBHOOK polled the api.
// This avoids tight culling with the other endpoints, and keeps all the caching and complexity inside the notifications package.

// Webhook is the body of the incoming request to be parsed and registered.
type Webhook struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout"`
	Field   string `json:"field"`
	Country string `json:"country"`
	Trigger string `json:"trigger"`
}

// Invoke a webhook by figuring out what it's looking for and getting it.
func (w *Webhook) Invoke() error {
	var body interface{}

	if w.Field == FieldConfirmed {
		confirmed, err := corona.GetLatestCases(w.Country)
		if err != nil {
			return err
		}
		body = &confirmed
	} else { // w.Field == FieldStringency
		stringency, err := corona.GetLatestStringency(w.Country)
		if err != nil {
			return err
		}
		body = &stringency
	}

	// Create a post request where the body is the data associated with the Webhooks field.
	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(body)
	req, err := http.NewRequest(http.MethodPost, w.URL, payload)
	if err != nil {
		return err
	}

	// Send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()

	if !corona.StatusIs2XX(res.StatusCode) {
		return &corona.ServerError{Err: "Remote responded with non 2XX code", StatusCode: res.StatusCode}
	}

	return nil
}

// Loop over webhooks and do stuff
func InvokeLoop(fs *firestore.Client, wg *sync.WaitGroup) {
	docref := fs.Collection(WebhookCollection).Doc("rqK-Q6dMa81bQF0UkdY4b8jQrpeb1BG6ahH5tXjctD4")

	docsnap, _ := docref.Get(context.Background())

	var data Webhook
	_ = docsnap.DataTo(&data)

	_ = data.Invoke()

	wg.Done()
}
