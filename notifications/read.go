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

// Just fetch the document by document id and return it.

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

// Just fetch the collection and return all documents in it.

func NewReadAllHandler(fs *firestore.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
