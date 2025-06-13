package cache

import (
	"context"
	"github.com/amundfpl/Assignment-2/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSetAndGetCountryCache(t *testing.T) {
	ctx := context.Background()
	key := "TEST_NO"
	data := utils.CountryInfoResponse{
		Name: struct {
			Common string `json:"common"`
		}{Common: "Norway"},
		Capital:    []string{"Oslo"},
		Latlng:     []float64{59.91, 10.75},
		Population: 5000000,
		Area:       385207.0,
		Currencies: map[string]utils.CurrencyDetails{
			"NOK": {Name: "Norwegian krone", Symbol: "kr"},
		},
	}

	err := SaveCountryInfoToCache(ctx, key, data)
	assert.NoError(t, err)

	cached, err := GetCachedCountryInfo(ctx, key, 1*time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, data.Name.Common, cached.Name.Common)
	assert.Equal(t, data.Capital, cached.Capital)
}

func TestSetAndGetWeatherCache(t *testing.T) {
	ctx := context.Background()
	key := WeatherCacheKey(59.91, 10.75)
	data := utils.WeatherData{
		Temperature:   5.5,
		Precipitation: 0.8,
	}

	err := SaveWeatherToCache(ctx, key, data)
	assert.NoError(t, err)

	cached, err := GetCachedWeather(ctx, key, 1*time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, data.Temperature, cached.Temperature)
	assert.Equal(t, data.Precipitation, cached.Precipitation)
}

func TestSetAndGetCurrencyCache(t *testing.T) {
	ctx := context.Background()
	key := CurrencyCacheKey("NOK", []string{"USD", "EUR"})
	rates := map[string]float64{
		"USD": 0.10,
		"EUR": 0.09,
	}

	err := SaveCurrencyRatesToCache(ctx, key, rates)
	assert.NoError(t, err)

	cached, err := GetCachedCurrencyRates(ctx, key, 1*time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, rates["USD"], cached["USD"])
	assert.Equal(t, rates["EUR"], cached["EUR"])
}
