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

// The error handling in this file is basically, if there is an error, print it to stdout, and move on.
// There are very few cases where we can do much more than that.

var (
	LastConfirmed  corona.CountryResponse
	LastStringency corona.PolicyResponse
)

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
// Returns if the invocation resulted in new data.
func (w *Webhook) Invoke(fs *firestore.Client, id string) (bool, string, error) {
	var body interface{}
	changed := false

	// Get whatever info the webhook is interested in
	if w.Field == FieldConfirmed {
		confirmed, err := corona.GetLatestCases(w.Country)
		if err != nil {
			return false, "", err
		}
		body = &confirmed
		if confirmed != LastConfirmed {
			changed = true
			LastConfirmed = confirmed
		}
	} else { // w.Field == FieldStringency
		stringency, err := corona.GetLatestStringency(w.Country)
		if err != nil {
			return false, "", err
		}
		body = &stringency
		if stringency != LastStringency {
			changed = true
			LastStringency = stringency
		}
	}

	// Create a post request where the body is the data associated with the Webhooks field.
	payload := new(bytes.Buffer)
	_ = json.NewEncoder(payload).Encode(body)
	req, err := http.NewRequest(http.MethodPost, w.URL, payload)
	if err != nil {
		return false, "", err
	}

	// Send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, "", err
	}
	res.Body.Close()

	if !corona.StatusIs2XX(res.StatusCode) {
		return false, "", &corona.ServerError{Err: "Remote responded with non 2XX code", StatusCode: res.StatusCode}
	}

	// Update the LastTriggered field of the webhook to now
	_, err = fs.
		Collection(WebhookCollection).
		Doc(id).
		Update(context.Background(), []firestore.Update{{Path: "LastTriggered", Value: time.Now()}})
	if err != nil {
		return false, "", err
	}

	return changed, w.Field, nil
}

// GetNextTimeout return the next timepoint where a webhook should be invoked.
// Returns the time of the next timeout, the id and webhook body of the next webhook to get timed out, and maybe an error.
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

// InvokeAllWithField invokes all the webhooks with the specified field.
func InvokeAllWithField(fs *firestore.Client, field, id string) error {
	// Get refs to all the documents in the collection
	docrefs, err := fs.Collection(WebhookCollection).DocumentRefs(context.Background()).GetAll()
	if err != nil {
		return err
	}

	// For all the docrefs, get their document snapshot, and parse them into webhooks, and invoke them
	for _, docref := range docrefs {
		docsnap, err := docref.Get(context.Background())
		if err != nil {
			return err
		}

		// Parse the document into a webhook
		var webhook Webhook
		err = docsnap.DataTo(&webhook)
		if err != nil {
			return err
		}

		// Skip webhooks with wrong trigger
		if webhook.Trigger != TriggerOnChange {
			continue
		}

		// Skip webhooks with different fields
		if webhook.Field != field {
			continue
		}

		// Skip the webhook that caused the refresh
		if docref.ID == id {
			continue
		}

		// Invoke the webhook, we can ignore any changes that results from this
		_, _, err = webhook.Invoke(fs, docref.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Loop over webhooks and invoke them if the should be invoked.
func InvokeLoop(fs *firestore.Client, wg *sync.WaitGroup) {
	for {
		next, id, webhook, err := GetNextTimeout(fs)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// Has the timeout already expired?
		if next.Before(time.Now()) {
			changed, field, err := webhook.Invoke(fs, id)
			if err != nil {
				log.Println(err.Error())
			}
			if changed {
				err = InvokeAllWithField(fs, field, id)
				if err != nil {
					log.Println(err.Error())
				}
			}
			continue
		}

		// Sleep until next timeout
		time.Sleep(time.Until(next))

		// Now invoke the webhook, since it has timed out by now
		changed, field, err := webhook.Invoke(fs, id)
		if err != nil {
			log.Println(err.Error())
		}
		if changed {
			err = InvokeAllWithField(fs, field, id)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	wg.Done()
}
