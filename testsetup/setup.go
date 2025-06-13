package testsetup

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
)

// InitTestFirebase sets up a Firebase app and Firestore client for integration testing.
// It loads service account credentials and connects to a dedicated test project.
func InitTestFirebase() {
	ctx := context.Background()

	// Load service account credentials path from utils
	credentialsPath := utils.DefaultCredentialsPath()

	// Initialize Firebase app using test project ID
	opts := option.WithCredentialsFile(credentialsPath)
	testApp, firebaseInitErr := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: utils.TestFirebaseProjectID,
	}, opts)
	if firebaseInitErr != nil {
		log.Fatalf(utils.ErrFirebaseAppInitializationFailed, firebaseInitErr)
	}

	// Create Firestore client from initialized app
	testFirestoreClient, firestoreClientErr := testApp.Firestore(ctx)
	if firestoreClientErr != nil {
		log.Fatalf(utils.ErrFirestoreClientInitializationFailed, firestoreClientErr)
	}

	// Assign the test Firestore client to the database layer
	db.SetClient(testFirestoreClient)
}
