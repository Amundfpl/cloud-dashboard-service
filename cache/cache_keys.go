package cache

import (
	"fmt"
	"github.com/amundfpl/Assignment-2/utils"
	"sort"
	"strings"
)

// WeatherCacheKey generates a unique cache key for a weather lookup based on latitude and longitude.
// It uses a formatted string defined in utils.WeatherCacheKeyFormat.
func WeatherCacheKey(lat, lon float64) string {
	return fmt.Sprintf(utils.WeatherCacheKeyFormat, lat, lon)
}

// CurrencyCacheKey generates a deterministic cache key for currency conversion based on
// a base currency and a list of target currencies.
// The target currencies are sorted to ensure the key is consistent regardless of input order.
func CurrencyCacheKey(base string, targets []string) string {
	sort.Strings(targets) // ensure deterministic key
	return base + utils.CacheKeySeparator + strings.Join(targets, utils.CacheKeySeparator)
}

// CountryCacheKey normalizes an ISO country code by trimming whitespace and converting to uppercase.
// This ensures consistency when storing or looking up country data in the cache.
func CountryCacheKey(iso string) string {
	return strings.ToUpper(strings.TrimSpace(iso))
}
