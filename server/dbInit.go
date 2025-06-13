package server

import (
	"context"
	"fmt"
	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
)

// DatabaseInitialization sets up the Firebase and Firestore clients.
// It must be called before any database operations are performed.
// Returns an error if either initialization step fails.
func DatabaseInitialization() error {
	ctx := context.Background()

	// Initialize Firebase app using credentials from the utils package
	if firebaseErr := db.InitFirebase(ctx, utils.CredentialsFirebaseKey); firebaseErr != nil {
		return fmt.Errorf(utils.ErrFirebaseInit, firebaseErr)
	}

	// Initialize Firestore client based on the Firebase app
	if firestoreErr := db.InitFirestore(ctx); firestoreErr != nil {
		return fmt.Errorf(utils.ErrFirestoreInit, firestoreErr)
	}

	return nil
}
