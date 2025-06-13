package db

import (
	"context"
	"fmt"
	"github.com/amundfpl/Assignment-2/utils"
	"google.golang.org/api/iterator"
	"log"
)

// InitFirestore initializes the global Firestore client using the Firebase application instance.
// This should be called once on startup before accessing Firestore.
func InitFirestore(ctx context.Context) error {
	client, initErr := FirebaseApp.Firestore(ctx)
	if initErr != nil {
		return fmt.Errorf(utils.ErrInitFirestoreClient, initErr)
	}
	firestoreClient = client
	return nil
}

// SaveDashboardConfig stores a new dashboard configuration in Firestore.
// It returns the generated document ID or an error.
func SaveDashboardConfig(ctx context.Context, config utils.DashboardConfig) (string, error) {
	docRef, _, saveErr := firestoreClient.Collection(utils.DashboardCollection).Add(ctx, config)
	if saveErr != nil {
		return "", saveErr
	}
	log.Println(utils.MsgDashboardSaved, docRef.ID)
	return docRef.ID, nil
}

// GetDashboardConfigByID retrieves a dashboard configuration by its document ID.
// It returns a pointer to the config or an error if not found or failed to decode.
func GetDashboardConfigByID(ctx context.Context, id string) (*utils.DashboardConfig, error) {
	docSnap, getErr := firestoreClient.Collection(utils.DashboardCollection).Doc(id).Get(ctx)
	if getErr != nil {
		return nil, getErr
	}

	var config utils.DashboardConfig
	decodeErr := docSnap.DataTo(&config)
	if decodeErr != nil {
		return nil, decodeErr
	}

	config.ID = docSnap.Ref.ID // Attach the document ID to the config object
	return &config, nil
}

// GetAllDashboardConfigs retrieves all dashboard configurations from Firestore.
// Returns a slice of configs or an error if fetching or decoding fails.
func GetAllDashboardConfigs(ctx context.Context) ([]utils.DashboardConfig, error) {
	iter := firestoreClient.Collection(utils.DashboardCollection).Documents(ctx)
	var configs []utils.DashboardConfig

	for {
		docSnap, nextErr := iter.Next()
		if nextErr == iterator.Done {
			break
		}
		if nextErr != nil {
			return nil, nextErr
		}

		var config utils.DashboardConfig
		decodeErr := docSnap.DataTo(&config)
		if decodeErr != nil {
			return nil, decodeErr
		}
		config.ID = docSnap.Ref.ID
		configs = append(configs, config)
	}
	return configs, nil
}

// UpdateDashboardConfig overwrites an existing dashboard configuration in Firestore
// based on its ID. Returns an error if the operation fails.
func UpdateDashboardConfig(ctx context.Context, config utils.DashboardConfig) error {
	_, updateErr := firestoreClient.Collection(utils.DashboardCollection).Doc(config.ID).Set(ctx, config)
	return updateErr
}

// DeleteDashboardConfig removes a dashboard configuration document by its ID.
// Returns an error if the delete operation fails.
func DeleteDashboardConfig(ctx context.Context, id string) error {
	_, deleteErr := firestoreClient.Collection(utils.DashboardCollection).Doc(id).Delete(ctx)
	return deleteErr
}

// PingFirestore checks Firestore connectivity by attempting to fetch one collection.
// Returns nil if successful, or an error otherwise.
func PingFirestore(ctx context.Context) error {
	_, pingErr := firestoreClient.Collections(ctx).Next()
	return pingErr // nil = Firestore is reachable
}
