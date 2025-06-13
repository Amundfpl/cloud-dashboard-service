package services

import (
	"context"
	"encoding/json"
	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterDashboardConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"name": map[string]interface{}{"common": "Mockistan"}},
		})
	}))
	defer server.Close()
	utils.RESTCountriesAPI = server.URL

	payload := []byte(`{
		"isoCode": "MOCK",
		"features": {"temperature": true, "capital": true}
	}`)

	resp, err := RegisterDashboardConfig(payload)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp["id"] == "" || resp["lastChange"] == "" {
		t.Fatal("Expected id and lastChange in response")
	}
}

func TestUpdateDashboardConfig(t *testing.T) {
	ctx := context.Background()
	doc := utils.DashboardConfig{
		Country:  "OldLand",
		ISOCode:  "OLD",
		Features: utils.FeatureConfig{Capital: true},
	}
	id, _ := db.SaveDashboardConfig(ctx, doc)

	body := []byte(`{
		"country": "NewLand",
		"isoCode": "NEW",
		"features": {"capital": true, "temperature": true}
	}`)

	resp, err := UpdateDashboardConfig(ctx, id, body)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp["id"] != id {
		t.Errorf("Expected ID %s, got %s", id, resp["id"])
	}
}

func TestPatchDashboardConfig(t *testing.T) {
	ctx := context.Background()
	cfg := utils.DashboardConfig{
		Country:  "PatchLand",
		ISOCode:  "PCH",
		Features: utils.FeatureConfig{Capital: false, Area: false},
	}
	id, _ := db.SaveDashboardConfig(ctx, cfg)

	patch := map[string]interface{}{
		"features": map[string]interface{}{
			"capital": true,
			"area":    true,
		},
	}
	resp, err := PatchDashboardConfig(ctx, id, patch)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if resp["id"] != id {
		t.Errorf("Expected ID %s, got %s", id, resp["id"])
	}
}

func TestDeleteRegistrationByID(t *testing.T) {
	ctx := context.Background()
	cfg := utils.DashboardConfig{
		Country:  "DelLand",
		ISOCode:  "DEL",
		Features: utils.FeatureConfig{},
	}
	id, _ := db.SaveDashboardConfig(ctx, cfg)

	err := DeleteRegistrationByID(ctx, id)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	_, err = db.GetDashboardConfigByID(ctx, id)
	if err == nil {
		t.Fatal("Expected error getting deleted config, got none")
	}
}
