package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amundfpl/Assignment-2/cache"
	"github.com/amundfpl/Assignment-2/httpclient"
	"github.com/amundfpl/Assignment-2/utils"
)

func TestGetEnrichedDashboards(t *testing.T) {
	// Unique ISO to avoid test collisions
	isoCode := fmt.Sprintf("MOCK_%d", time.Now().UnixNano())
	currency := "FOO"
	lat, lon := 42.42, 24.24

	// Create and cache country info
	countryInfo := utils.CountryInfoResponse{
		Capital:    []string{"TestCity"},
		Latlng:     []float64{lat, lon},
		Population: 123456,
		Area:       7890,
		Currencies: map[string]utils.CurrencyDetails{currency: {Name: "FakeCoin"}},
	}
	_ = cache.SaveCountryInfoToCache(context.Background(), isoCode, countryInfo)

	// Cache weather
	weatherKey := cache.WeatherCacheKey(lat, lon)
	_ = cache.SaveWeatherToCache(context.Background(), weatherKey, utils.WeatherData{
		Temperature:   -3.5,
		Precipitation: 1.2,
	})

	// Cache currency
	currencyKey := cache.CurrencyCacheKey(currency, []string{"USD", "EUR"})
	_ = cache.SaveCurrencyRatesToCache(context.Background(), currencyKey, map[string]float64{
		"USD": 1.23,
		"EUR": 0.95,
	})

	// Build test dashboard config directly (bypassing Firestore)
	config := utils.DashboardConfig{
		Country: "Mockistan",
		ISOCode: isoCode,
		Features: utils.FeatureConfig{
			Capital:          true,
			Coordinates:      true,
			Population:       true,
			Area:             true,
			Temperature:      true,
			Precipitation:    true,
			TargetCurrencies: []string{"USD", "EUR"},
		},
	}

	client := httpclient.NewClient()
	resp := utils.DashboardResponse{
		Country: config.Country,
		ISOCode: config.ISOCode,
	}

	cInfo, err := enrichCountryData(client, config, &resp)
	if err != nil {
		t.Fatalf("enrichCountryData failed: %v", err)
	}
	if resp.Capital != "TestCity" {
		t.Errorf("Expected capital TestCity, got %s", resp.Capital)
	}

	err = enrichWeatherData(client, config, &resp)
	if err != nil {
		t.Fatalf("enrichWeatherData failed: %v", err)
	}
	if resp.Temperature != -3.5 {
		t.Errorf("Expected temperature -3.5, got %f", resp.Temperature)
	}

	err = enrichCurrencyData(client, config, cInfo, &resp)
	if err != nil {
		t.Fatalf("enrichCurrencyData failed: %v", err)
	}
	if resp.ExchangeRates["USD"] != 1.23 {
		t.Errorf("Expected USD rate 1.23, got %f", resp.ExchangeRates["USD"])
	}
	if resp.ExchangeRates["EUR"] != 0.95 {
		t.Errorf("Expected EUR rate 0.95, got %f", resp.ExchangeRates["EUR"])
	}
}

// Test syncCountryFields independently
func TestSyncCountryFields(t *testing.T) {
	resp := &utils.DashboardResponse{}
	cfg := utils.DashboardConfig{
		Features: utils.FeatureConfig{
			Capital:     true,
			Coordinates: true,
			Population:  true,
			Area:        true,
		},
	}
	info := utils.CountryInfoResponse{
		Capital:    []string{"CapCity"},
		Latlng:     []float64{60.0, 5.0},
		Population: 1000000,
		Area:       300.5,
	}
	syncCountryFields(cfg, resp, info)

	if resp.Capital != "CapCity" {
		t.Errorf("Expected CapCity, got %s", resp.Capital)
	}
	if resp.Latitude != 60.0 || resp.Longitude != 5.0 {
		t.Errorf("Expected coordinates 60.0/5.0, got %f/%f", resp.Latitude, resp.Longitude)
	}
	if resp.Population != 1000000 {
		t.Errorf("Expected population 1000000, got %d", resp.Population)
	}
	if resp.Area != 300.5 {
		t.Errorf("Expected area 300.5, got %f", resp.Area)
	}
}

// Force fetch path of enrichWeatherData (cache miss)
func TestEnrichWeatherData_WithoutCache(t *testing.T) {
	cfg := utils.DashboardConfig{
		Features: utils.FeatureConfig{
			Temperature:   true,
			Precipitation: true,
		},
	}
	// Use random coordinates to avoid hitting old cache
	lat := 51.5 + float64(time.Now().UnixNano()%1000)/10000.0
	lon := -0.1 + float64(time.Now().UnixNano()%1000)/10000.0
	resp := &utils.DashboardResponse{Latitude: lat, Longitude: lon}
	client := httpclient.NewClient()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]float64{
				"temperature_2m": 5.5,
				"precipitation":  1.1,
			},
		})
	}))
	defer weatherServer.Close()
	utils.OpenMeteoAPI = weatherServer.URL

	err := enrichWeatherData(client, cfg, resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Temperature != 5.5 || resp.Precipitation != 1.1 {
		t.Errorf("Expected 5.5°C and 1.1mm, got %f°C and %fmm", resp.Temperature, resp.Precipitation)
	}
}

// Force fetch path of enrichCurrencyData (cache miss)
func TestEnrichCurrencyData_WithoutCache(t *testing.T) {
	cfg := utils.DashboardConfig{
		Features: utils.FeatureConfig{
			TargetCurrencies: []string{"USD"},
		},
	}
	resp := &utils.DashboardResponse{}
	client := httpclient.NewClient()

	// Randomize base currency to avoid hitting cached values
	base := fmt.Sprintf("TEST_%d", time.Now().UnixNano())

	currencyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"rates": map[string]float64{
				"USD": 1.33,
			},
		})
	}))
	defer currencyServer.Close()
	utils.CurrencyAPI = currencyServer.URL

	countryInfo := utils.CountryInfoResponse{
		Currencies: map[string]utils.CurrencyDetails{
			base: {Name: "TestCoin"},
		},
	}

	err := enrichCurrencyData(client, cfg, countryInfo, resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.ExchangeRates["USD"] != 1.33 {
		t.Errorf("Expected USD rate 1.33, got %f", resp.ExchangeRates["USD"])
	}
}
