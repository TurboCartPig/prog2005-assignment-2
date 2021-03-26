package notifications

import (
	"cloud.google.com/go/firestore"
	"net/http"
)

func NewDeleteHandler(fs *firestore.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}
