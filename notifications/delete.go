package notifications

import (
	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

// NewDeleteHandler creates a HttpHandler that, given a webhook id, deletes the webhook from the database.
func NewDeleteHandler(fs *firestore.Client) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		// TODO: Validate the id somehow?

		// Get the document out of the collection
		docref := fs.Collection(WebhookCollection).Doc(id)
		if docref == nil {
			http.Error(rw, "Id not valid, the webhooks either does not exits, or it has already been deleted.", http.StatusBadRequest)
			return
		}

		// And delete it
		_, err := docref.Delete(r.Context())
		if err != nil {
			http.Error(rw, "Something went wrong when deleting the webhook.", http.StatusInternalServerError)
			return
		}

		log.Println("Deleting webhook with id:", id)
		rw.WriteHeader(http.StatusNoContent)
	}
}
