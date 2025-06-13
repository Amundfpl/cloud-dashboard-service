package server

import (
	"context"
	"fmt"
	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
	"log"
	"os"
)

// DatabaseInitialization sets up the Firebase and Firestore clients.
// It must be called before any database operations are performed.
// Returns an error if either initialization step fails.
func DatabaseInitialization() error {
	ctx := context.Background()

	// Initialize Firebase app using credentials from the utils package
	firebasePath := os.Getenv(utils.GOOGLE_APPLICATION_CREDENTIALS)
	if firebasePath == "" {
		log.Println(utils.ErrGOOGLE_APPLICATION_CREDENTIALS_NotSet)
		firebasePath = utils.CredentialsFirebaseKey
	}

	if err := db.InitFirebase(ctx, firebasePath); err != nil {
		return fmt.Errorf(utils.ErrFirebaseInit, err)
	}

	// Initialize Firestore client based on the Firebase app
	if firestoreErr := db.InitFirestore(ctx); firestoreErr != nil {
		return fmt.Errorf(utils.ErrFirestoreInit, firestoreErr)
	}

	return nil
}
