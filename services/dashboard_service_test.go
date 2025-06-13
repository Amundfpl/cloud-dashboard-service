// services/dashboard_service_test.go
package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/httpclient"
	"github.com/amundfpl/Assignment-2/utils"
)

import (
	"github.com/amundfpl/Assignment-2/testsetup"
	"os"
)

func TestMain(m *testing.M) {
	testsetup.InitTestFirebase()
	clearTestCacheCollections()
	os.Exit(m.Run())
}

func clearTestCacheCollections() {
	ctx := context.Background()
	collections := []string{"country_cache", "weather_cache", "currency_cache"}
	for _, col := range collections {
		docs, _ := db.GetClient().Collection(col).Documents(ctx).GetAll()
		for _, doc := range docs {
			_, _ = doc.Ref.Delete(ctx)
		}
	}
}

func TestGetPopulatedDashboardByID(t *testing.T) {
	testID := "dash-test-service"
	ctx := context.Background()

	_, err := db.GetClient().Collection("dashboard_configs").Doc(testID).Set(ctx, map[string]interface{}{
		"country": "Mockland",
		"isoCode": "NO",
		"features": map[string]interface{}{
			"capital":          true,
			"coordinates":      true,
			"population":       true,
			"area":             true,
			"temperature":      true,
			"precipitation":    true,
			"targetCurrencies": []string{"USD", "EUR"},
		},
	})
	if err != nil {
		t.Fatalf(" Failed to seed test config: %v", err)
	}

	countryStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{
			"capital": ["Mockville"],
			"latlng": [60.0, 10.0],
			"population": 100000,
			"area": 5432.1,
			"currencies": {"NOK": {"name": "Krone"}}
		}]`))
	}))
	defer countryStub.Close()

	weatherStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"current": {
				"temperature_2m": -5.2,
				"precipitation": 2.1
			}
		}`))
	}))
	defer weatherStub.Close()

	currencyStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"rates": {"USD": 1.1, "EUR": 0.9}
		}`))
	}))
	defer currencyStub.Close()

	// Override global API URLs
	originalCountryAPI := utils.RESTCountriesAPI
	originalWeatherAPI := utils.OpenMeteoAPI
	originalCurrencyAPI := utils.CurrencyAPI
	utils.RESTCountriesAPI = countryStub.URL
	utils.OpenMeteoAPI = weatherStub.URL
	utils.CurrencyAPI = currencyStub.URL
	defer func() {
		utils.RESTCountriesAPI = originalCountryAPI
		utils.OpenMeteoAPI = originalWeatherAPI
		utils.CurrencyAPI = originalCurrencyAPI
	}()

	resp, err := GetPopulatedDashboardByID(testID)
	if err != nil {
		t.Fatalf("Failed to get populated dashboard: %v", err)
	}

	if resp.Country != "Mockland" || resp.Features.Capital != "Mockville" {
		t.Errorf("Invalid data: %+v", resp)
	}
	if resp.Features.TargetCurrencies["USD"] != 1.1 {
		t.Errorf("Expected USD=1.1, got: %f", resp.Features.TargetCurrencies["USD"])
	}
}

func TestFetchCountryInfo(t *testing.T) {
	client := httpclient.NewClient()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{
			"capital": ["Mockville"],
			"latlng": [50.0, 20.0],
			"population": 100000,
			"area": 9876.5,
			"currencies": {"XYZ": {"name": "Mockcoin"}}
		}]`))
	}))
	defer server.Close()

	utils.RESTCountriesAPI = server.URL
	info, err := fetchCountryInfo(client, "XYZ")
	if err != nil {
		t.Fatalf("Error fetching country info: %v", err)
	}
	if info.Population != 100000 {
		t.Errorf("Expected population 100000, got: %d", info.Population)
	}
}

func TestFetchWeather(t *testing.T) {
	client := httpclient.NewClient()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"current": {
				"temperature_2m": 10.5,
				"precipitation": 0.3
			}
		}`))
	}))
	defer server.Close()

	utils.OpenMeteoAPI = server.URL
	weather, err := fetchWeather(client, 60.0, 10.0)
	if err != nil {
		t.Fatalf("Error fetching weather: %v", err)
	}
	if weather.Temperature != 10.5 || weather.Precipitation != 0.3 {
		t.Errorf("Unexpected weather values: %+v", weather)
	}
}

func TestFetchCurrencyRates(t *testing.T) {
	client := httpclient.NewClient()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"rates": {
			"EUR": 0.9,
			"NOK": 10.5
  }
}
`))
	}))
	defer server.Close()

	utils.CurrencyAPI = server.URL
	currencies := map[string]utils.CurrencyDetails{"NOK": {Name: "Krone"}}
	rates, err := fetchCurrencyRates(client, currencies, []string{"USD", "EUR"})
	if err != nil {
		t.Fatalf("Error fetching currency rates: %v", err)
	}
	if rates["USD"] != 1.23 {
		t.Errorf("Expected USD=1.23, got: %f", rates["USD"])
	}
	if rates["EUR"] != 0.98 {
		t.Errorf("Expected EUR=0.98, got: %f", rates["EUR"])
	}
}
