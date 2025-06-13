package db

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/amundfpl/Assignment-2/utils"
	"google.golang.org/api/option"
)

var (
	FirebaseApp     *firebase.App
	firestoreClient *firestore.Client
)

// SetClient sets the global Firestore client (typically used in tests).
func SetClient(client *firestore.Client) {
	firestoreClient = client
}

// GetClient returns the global Firestore client.
func GetClient() *firestore.Client {
	return firestoreClient
}

// InitFirebase initializes the Firebase app and Firestore client.
// Uses credentials from the provided path and the project ID defined in utils.
func InitFirebase(ctx context.Context, credentialsFile string) error {
	opt := option.WithCredentialsFile(credentialsFile)

	config := &firebase.Config{
		ProjectID: utils.FirebaseProjectID,
	}

	app, appErr := firebase.NewApp(ctx, config, opt)
	if appErr != nil {
		return fmt.Errorf(utils.ErrInitFirebaseApp, appErr)
	}

	client, clientErr := app.Firestore(ctx)
	if clientErr != nil {
		return fmt.Errorf(utils.ErrInitFirestoreClient, clientErr)
	}

	FirebaseApp = app
	firestoreClient = client
	return nil
}

// FirestoreClient exposes the Firestore client or panics if it's uninitialized.
// Used internally to ensure valid access.
func FirestoreClient() *firestore.Client {
	if firestoreClient == nil {
		panic(utils.ErrFirestoreNotInitialized)
	}
	return firestoreClient
}

// CloseFirestore releases the Firestore client connection, if initialized.
func CloseFirestore() error {
	if firestoreClient != nil {
		return firestoreClient.Close()
	}
	return nil
}

// IsFirestoreInitialized checks if the Firestore client is ready.
// Useful in test environments to avoid panics.
func IsFirestoreInitialized() bool {
	return firestoreClient != nil
}
