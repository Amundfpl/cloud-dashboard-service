package handlers

import (
	"context"
	"encoding/json"
	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/services"
	"github.com/amundfpl/Assignment-2/testsetup"
	"github.com/amundfpl/Assignment-2/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestMain initializes the Firebase test client once before all tests run
func TestMain(m *testing.M) {
	testsetup.InitTestFirebase()
	m.Run()
}

// TestHandleGetPopulatedDashboard_RealService verifies that the GET /dashboards/{id} handler returns enriched data
func TestHandleGetPopulatedDashboard_RealService(t *testing.T) {
	testID := "dashboard-test-123"
	client := db.GetClient()
	ctx := context.Background()

	_, err := client.Collection("dashboard_configs").Doc(testID).Set(ctx, map[string]interface{}{
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
		t.Fatalf("Failed to set test config: %v", err)
	}

	// --- Stub external APIs ---

	countryStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `[{
			"capital": ["Mockville"],
			"latlng": [59.9139, 10.7522],
			"population": 123456,
			"area": 6543.21,
			"currencies": {"NOK": {"name": "Krone"}}
		}]`
		w.Write([]byte(resp))
	}))
	defer countryStub.Close()

	weatherStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"current": {
				"temperature_2m": -2.0,
				"precipitation": 1.5
			}
		}`))
	}))
	defer weatherStub.Close()

	currencyStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"rates": {
				"USD": 1.23,
				"EUR": 1.05
			}
		}`))
	}))
	defer currencyStub.Close()

	// Temporarily override API base URLs
	originalCountries := utils.RESTCountriesAPI
	originalWeather := utils.OpenMeteoAPI
	originalCurrency := utils.CurrencyAPI

	utils.RESTCountriesAPI = countryStub.URL
	utils.OpenMeteoAPI = weatherStub.URL
	utils.CurrencyAPI = currencyStub.URL

	defer func() {
		utils.RESTCountriesAPI = originalCountries
		utils.OpenMeteoAPI = originalWeather
		utils.CurrencyAPI = originalCurrency
	}()

	// --- Execute handler under test ---
	service := services.RealDashboardService{}
	getHandler := NewDashboardHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/dashboards/"+testID, nil)
	rr := httptest.NewRecorder()

	getHandler(rr, req)

	// --- Assert output ---
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response utils.PopulatedDashboardResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.Country != "Mockland" {
		t.Errorf("Expected 'Mockland', got: %s", response.Country)
	}
	if response.Features.Capital != "Mockville" {
		t.Errorf("Expected capital to be 'Mockville', got: %s", response.Features.Capital)
	}
	if response.Features.TargetCurrencies["USD"] != 1.23 {
		t.Errorf("Expected USD rate to be 1.23, got: %f", response.Features.TargetCurrencies["USD"])
	}
}
