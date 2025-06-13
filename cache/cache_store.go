package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
)

// cacheEntry is a generic cache structure that wraps the data with a timestamp.
// This helps in tracking when the data was cached and checking for expiration.
type cacheEntry[T any] struct {
	Timestamp time.Time
	Data      T
}

// isCacheExpired returns true if the given timestamp is older than the allowed maximum age.
func isCacheExpired(timestamp time.Time, maxAge time.Duration) bool {
	return time.Since(timestamp) > maxAge
}

// setCache stores a generic value into Firestore with a timestamp in the specified collection under the given document ID.
func setCache[T any](ctx context.Context, collection, docID string, data T) error {
	_, err := db.FirestoreClient().Collection(collection).Doc(docID).Set(ctx, map[string]interface{}{
		utils.FieldData:      data,
		utils.TimestampField: time.Now(),
	})
	return err
}

// getCache retrieves a value from Firestore, checks if it's expired, and returns the typed data.
// It returns an error if the document is missing, decoding fails, or the data is too old.
func getCache[T any](ctx context.Context, collection, docID string, maxAge time.Duration) (*T, error) {
	doc, err := db.FirestoreClient().Collection(collection).Doc(docID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrCacheMiss, docID, err)
	}

	var entry cacheEntry[T]
	if err := doc.DataTo(&entry); err != nil {
		return nil, fmt.Errorf(utils.ErrCacheDecode, docID, err)
	}

	if isCacheExpired(entry.Timestamp, maxAge) {
		return nil, errors.New(utils.ErrCacheExpired)
	}

	return &entry.Data, nil
}

// --- Country Cache ---

// GetCachedCountryInfo retrieves cached country data for a given ISO code if it is not expired.
func GetCachedCountryInfo(ctx context.Context, iso string, maxAge time.Duration) (*utils.CountryInfoResponse, error) {
	key := CountryCacheKey(iso)
	return getCache[utils.CountryInfoResponse](ctx, utils.CountryCacheCollection, key, maxAge)
}

// SaveCountryInfoToCache stores country data in the cache for the given ISO code.
func SaveCountryInfoToCache(ctx context.Context, iso string, data utils.CountryInfoResponse) error {
	key := CountryCacheKey(iso)
	return setCache(ctx, utils.CountryCacheCollection, key, data)
}

// --- Weather Cache ---

// GetCachedWeather retrieves cached weather data by key if it is not expired.
func GetCachedWeather(ctx context.Context, key string, maxAge time.Duration) (*utils.WeatherData, error) {
	return getCache[utils.WeatherData](ctx, utils.WeatherCacheCollection, key, maxAge)
}

// SaveWeatherToCache stores weather data in the cache under the given key.
func SaveWeatherToCache(ctx context.Context, key string, data utils.WeatherData) error {
	return setCache(ctx, utils.WeatherCacheCollection, key, data)
}

// --- Currency Cache ---

// GetCachedCurrencyRates retrieves cached currency exchange rates if available and not expired.
// Unlike other types, this uses a manual struct instead of the generic cacheEntry due to map typing.
func GetCachedCurrencyRates(ctx context.Context, key string, maxAge time.Duration) (map[string]float64, error) {
	doc, err := db.FirestoreClient().Collection(utils.CurrencyCacheCollection).Doc(key).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf(utils.ErrCacheMissCurrency, key, err)
	}

	// Manually define struct since map[string]float64 doesn't work with generics directly
	var entry struct {
		Timestamp time.Time
		Data      map[string]float64
	}
	if err := doc.DataTo(&entry); err != nil {
		return nil, fmt.Errorf(utils.ErrCacheDecodeCurrency, key, err)
	}

	if isCacheExpired(entry.Timestamp, maxAge) {
		return nil, errors.New(utils.ErrCacheExpiredCurrency)
	}

	return entry.Data, nil
}

// SaveCurrencyRatesToCache stores currency exchange rates under the given key in the cache.
func SaveCurrencyRatesToCache(ctx context.Context, key string, rates map[string]float64) error {
	return setCache(ctx, utils.CurrencyCacheCollection, key, rates)
}
