package notifications

import (
	"assignment-2/corona"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// There are two types of webhooks
// 1. ON_TIMEOUT, which is easy
// 2. ON_CHANGE, but how to know if something changed?
// I will make two queues with ON_CHANGE webhooks. One for "stringency" and one for "confirmed".
// In each queue I wait for the first timeout.
// The I respond to all the webhooks if the data has changed from last time a WEBHOOK polled the api.
// This avoids tight culling with the other endpoints, and keeps all the caching and complexity inside the notifications package.

// Webhook is the body of the any request involving a webhook.
type Webhook struct {
	URL           string    `json:"url"`
	Timeout       int       `json:"timeout"`
	Field         string    `json:"field"`
	Country       string    `json:"country"`
	Trigger       string    `json:"trigger"`
	LastTriggered time.Time `json:"last_triggered"`
}

// Invoke a webhook by figuring out what it's looking for and getting it.
func (w *Webhook) Invoke(fs *firestore.Client, id string) error {
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

	// Update the LastTriggered field of the webhook to now
	_, err = fs.
		Collection(WebhookCollection).
		Doc(id).
		Update(context.Background(), []firestore.Update{{Path: "LastTriggered", Value: time.Now()}})
	if err != nil {
		return err
	}

	return nil
}

// GetNextTimeout return the next timepoint where a webhook should be invoked.
func GetNextTimeout(fs *firestore.Client) (time.Time, string, *Webhook, error) {
	// Get refs to all the documents in the collection
	docrefs, err := fs.Collection(WebhookCollection).DocumentRefs(context.Background()).GetAll()
	if err != nil {
		return time.Now(), "", nil, err
	}

	// Initialize target to after any webhook timeout point
	target := time.Now().AddDate(100, 0, 0) // nolint:gomnd // Just some random large duration
	// Save the latest
	var targetWebhook *Webhook
	var targetID string
	for _, docref := range docrefs {
		docsnap, err := docref.Get(context.Background())
		if err != nil {
			continue // Just ignore any documents that can't be fetched, either because they don't have snapshots, or some other error
		}

		// Parse the document into webhook
		var data Webhook
		err = docsnap.DataTo(&data)
		if err != nil {
			continue
		}

		// Get the timeout as a duration
		timeoutdur := time.Duration(data.Timeout) * time.Second

		// Is the timeout point of the webhook more recent than the most
		// recently recorded
		if current := data.LastTriggered.Add(timeoutdur); current.Before(target) {
			target = current
			targetWebhook = &data
			targetID = docref.ID
		}
	}

	return target, targetID, targetWebhook, nil
}

// Loop over webhooks and do stuff
func InvokeLoop(fs *firestore.Client, wg *sync.WaitGroup) {
	for {
		next, id, webhook, err := GetNextTimeout(fs)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// Has the timeout already expired?
		if next.Before(time.Now()) {
			err = webhook.Invoke(fs, id)
			if err != nil {
				log.Println(err.Error())
			}
			continue
		}

		// Sleep until next timeout
		time.Sleep(time.Until(next))

		// Now invoke the webhook, since it has timed out by now
		err = webhook.Invoke(fs, id)
		if err != nil {
			log.Println(err.Error())
		}
	}

	wg.Done()
}
