// cache/cache_keys_test.go
package cache

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherCacheKey(t *testing.T) {
	key := WeatherCacheKey(59.91, 10.75)
	assert.Equal(t, "59.9_10.8", key) // Matches fmt.Sprintf("%.1f_%.1f", ...)
}

func TestCurrencyCacheKey_Deterministic(t *testing.T) {
	targets := []string{"USD", "EUR", "SEK"}
	key1 := CurrencyCacheKey("NOK", targets)

	scrambled := []string{"SEK", "USD", "EUR"}
	key2 := CurrencyCacheKey("NOK", scrambled)

	assert.Equal(t, key1, key2)
	assert.True(t, strings.HasPrefix(key1, "NOK"))
}

func TestCountryCacheKey(t *testing.T) {
	key := CountryCacheKey("  no ")
	assert.Equal(t, "NO", key)
}
