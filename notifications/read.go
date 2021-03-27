package notifications

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"github.com/go-chi/chi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

// NewReadHandler creates a HttpHandler that reads one webhook from the database and returns it.
func NewReadHandler(fs *firestore.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		docref := fs.Collection(WebhookCollection).Doc(id)
		docsnap, err := docref.Get(r.Context())
		if status.Code(err) == codes.NotFound {
			http.Error(rw, "Invalid webhook id; No webhook registered by that id", http.StatusBadRequest)
			return
		} else if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Something went wrong trying to get the webhook", http.StatusInternalServerError)
			return
		}

		var data requestBody
		err = docsnap.DataTo(&data)
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Something went wrong trying to read the webhook", http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(rw).Encode(&data)
	}
}

// NewReadAllHandler creates a HttpHandler that reads all the webhooks from the database and returns them.
func NewReadAllHandler(fs *firestore.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Get refs to all the documents in the collection
		docrefs, err := fs.Collection(WebhookCollection).DocumentRefs(r.Context()).GetAll()
		if err != nil {
			log.Println(err.Error())
			http.Error(rw, "Something went wrong trying to get the webhooks", http.StatusInternalServerError)
			return
		}

		// For all the docrefs, get their document snapshot,
		// and read the data into a struct, then append that struct to `body`
		// NOTE: There might be docrefs for which there is no document, therefore we can not preallocate the body slice.
		body := make([]requestBody, 0)
		for _, docref := range docrefs {
			docsnap, err := docref.Get(r.Context())
			if status.Code(err) == codes.NotFound {
				http.Error(rw, "Invalid webhook id; No webhook registered by that id", http.StatusBadRequest)
				return
			} else if err != nil {
				log.Println(err.Error())
				http.Error(rw, "Something went wrong trying to get the webhook", http.StatusInternalServerError)
				return
			}

			var data requestBody
			err = docsnap.DataTo(&data)
			if err != nil {
				log.Println(err.Error())
				http.Error(rw, "Something went wrong trying to read the webhook", http.StatusInternalServerError)
				return
			}

			body = append(body, data)
		}

		_ = json.NewEncoder(rw).Encode(&body)
	}
}
