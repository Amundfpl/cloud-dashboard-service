package cache

import (
	"context"
	"log"
	"time"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
)

// purgeCacheCollection deletes all documents in a given Firestore collection that are older than the specified duration.
// - collection: the Firestore collection name
// - olderThan: documents with timestamps before (now - olderThan) will be purged
func purgeCacheCollection(ctx context.Context, collection string, olderThan time.Duration) error {
	// Calculate time threshold for stale documents
	threshold := time.Now().Add(-olderThan)

	// Build a Firestore query to find documents older than the threshold
	query := db.FirestoreClient().
		Collection(collection).
		Where(utils.TimestampField, utils.OperatorLessThan, threshold)

	// Execute the query
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return err // Return any query error
	}

	// Log how many documents will be purged
	log.Printf(utils.MsgPurgeSuccess, len(docs), collection)

	// Delete each matching document
	for _, doc := range docs {
		_, _ = doc.Ref.Delete(ctx) // Ignoring delete errors (optional: handle if needed)
	}
	return nil
}

// PurgeOldCountryCache purges outdated entries from the country cache based on its TTL setting.
func PurgeOldCountryCache(ctx context.Context) error {
	return purgeCacheCollection(ctx, utils.CountryCacheCollection, utils.CountryCacheTTL) //Purge old country cache every 24 hours
}

// PurgeOldWeatherCache purges outdated entries from the weather cache based on its TTL setting.
func PurgeOldWeatherCache(ctx context.Context) error {
	return purgeCacheCollection(ctx, utils.WeatherCacheCollection, utils.WeatherCacheTTL) // Purge old weather cache every 2 hour
}

// PurgeOldCurrencyCache purges outdated entries from the currency cache based on its TTL setting.
func PurgeOldCurrencyCache(ctx context.Context) error {
	return purgeCacheCollection(ctx, utils.CurrencyCacheCollection, utils.CurrencyCacheTTL) // Purge old currency cache every 12 hour
}
