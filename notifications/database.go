package notifications

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

// NewFirestoreClient creates and initializes a new firestore client.
func NewFirestoreClient() *firestore.Client {
	// Use GOOGLE_APPLICATION_CREDENTIALS env var to find the service account key
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
