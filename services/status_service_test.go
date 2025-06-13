package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
)

func TestGetSystemStatus(t *testing.T) {
	ctx := context.Background()

	// --- Countries API Stub ---
	countryStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/alpha/no" {
			t.Errorf("Unexpected countries API path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer countryStub.Close()
	utils.RESTCountriesAPI = countryStub.URL

	// --- Meteo API Stub ---
	meteoStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/forecast") {
			t.Errorf("Unexpected meteo API path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer meteoStub.Close()
	utils.OpenMeteoAPI = meteoStub.URL

	// --- Currency API Stub ---
	currencyStub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Accept both /NOK and NOK to make test agnostic to slash logic
		if !strings.HasSuffix(r.URL.Path, "NOK") {
			t.Errorf("Expected path ending in 'NOK', got '%s'", r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer currencyStub.Close()
	utils.CurrencyAPI = currencyStub.URL + "/" // Ensures base URL ends with slash

	// --- Ensure webhook exists ---
	_, err := db.GetClient().Collection("webhooks").Doc("test-webhook").Set(ctx, map[string]any{
		"url":     "http://example.com",
		"event":   "INVOKE",
		"country": "NO",
	})
	if err != nil {
		t.Fatalf("Failed to set test webhook: %v", err)
	}

	// --- Call system status ---
	status := GetSystemStatus(ctx)

	if status.CountriesAPI != http.StatusOK {
		t.Errorf("Expected CountriesAPI status 200, got %d", status.CountriesAPI)
	}
	if status.MeteoAPI != http.StatusOK {
		t.Errorf("Expected MeteoAPI status 200, got %d", status.MeteoAPI)
	}
	if status.CurrencyAPI != http.StatusOK {
		t.Errorf("Expected CurrencyAPI status 200, got %d", status.CurrencyAPI)
	}
	if status.NotificationDB != http.StatusOK {
		t.Errorf("Expected NotificationDB status 200, got %d", status.NotificationDB)
	}
	if status.Webhooks == 0 {
		t.Errorf("Expected non-zero webhook count")
	}
	if status.Version != "v1" {
		t.Errorf("Expected version v1, got %s", status.Version)
	}
	if status.UptimeInSeconds < 0 {
		t.Errorf("Expected positive uptime, got %d", status.UptimeInSeconds)
	}
}

func getSystemStatusWithMockedFirestore(ctx context.Context, pingFn func(context.Context) error) utils.StatusReport {
	checkFirestore := func(ctx context.Context) int {
		if err := pingFn(ctx); err != nil {
			return 0
		}
		return http.StatusOK
	}

	return utils.StatusReport{
		CountriesAPI:    http.StatusOK,
		MeteoAPI:        http.StatusOK,
		CurrencyAPI:     http.StatusOK,
		NotificationDB:  checkFirestore(ctx),
		Webhooks:        1,
		Version:         "v1",
		UptimeInSeconds: 1,
	}
}

func TestCheckFirestoreFailure(t *testing.T) {
	ctx := context.Background()
	status := getSystemStatusWithMockedFirestore(ctx, func(ctx context.Context) error {
		return context.DeadlineExceeded
	})

	if status.NotificationDB == http.StatusOK {
		t.Error("Expected Firestore check to fail, got 200")
	}
}
