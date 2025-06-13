package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Hjelpefunksjon for å opprette dashboard og returnere ID
func createTestDashboard(t *testing.T) string {
	body := `{
		"country": "TEST",
		"isoCode": "TST",
		"features": {
			"capital": true
		}
	}`
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/registrations", strings.NewReader(body))
	rr := httptest.NewRecorder()
	HandleRegisterDashboard(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("❌ Failed to create test dashboard, got %d", rr.Code)
	}

	var result map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &result)
	return result["id"]
}

func TestHandleRegisterDashboard_Success(t *testing.T) {
	id := createTestDashboard(t)
	if id == "" {
		t.Fatal("Expected dashboard ID in response, got empty")
	}
}

func TestHandleRegisterDashboard_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/registrations", strings.NewReader("invalid"))
	rr := httptest.NewRecorder()
	HandleRegisterDashboard(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid body, got %d", rr.Code)
	}
}

func TestGetAllRegistrations_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/registrations", nil)
	rr := httptest.NewRecorder()

	GetAllRegistrations(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "[") {
		t.Errorf("Expected JSON array in response, got: %s", rr.Body.String())
	}
}

func TestGetRegistrationByID_Success(t *testing.T) {
	id := createTestDashboard(t)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/registrations/"+id, nil)
	rr := httptest.NewRecorder()
	GetRegistrationByID(rr, req, id)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"TEST"`) {
		t.Errorf("Expected response to contain 'TEST', got: %s", rr.Body.String())
	}
}

func TestGetRegistrationByID_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/registrations/does-not-exist", nil)
	rr := httptest.NewRecorder()
	GetRegistrationByID(rr, req, "does-not-exist")

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", rr.Code)
	}
}

func TestUpdateDashboardRegistration_Success(t *testing.T) {
	id := createTestDashboard(t)

	body := `{"country": "SE", "features": {"area": true}}`
	req := httptest.NewRequest(http.MethodPut, "/dashboard/v1/registrations/"+id, strings.NewReader(body))
	rr := httptest.NewRecorder()

	UpdateDashboardRegistration(rr, req, id)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

func TestUpdateDashboardRegistration_BadBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/dashboard/v1/registrations/abc123", strings.NewReader("not-json"))
	rr := httptest.NewRecorder()

	UpdateDashboardRegistration(rr, req, "abc123")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestHeadCheckDashboard_Success(t *testing.T) {
	id := createTestDashboard(t)

	req := httptest.NewRequest(http.MethodHead, "/dashboard/v1/registrations/"+id, nil)
	rr := httptest.NewRecorder()

	HeadCheckDashboard(rr, req, id)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

func TestHeadCheckDashboard_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, "/dashboard/v1/registrations/unknown", nil)
	rr := httptest.NewRecorder()

	HeadCheckDashboard(rr, req, "unknown")

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", rr.Code)
	}
}

func TestPatchDashboardRegistration_Success(t *testing.T) {
	id := createTestDashboard(t)

	patch := map[string]interface{}{
		"country": "DK",
	}
	body, _ := json.Marshal(patch)

	req := httptest.NewRequest(http.MethodPatch, "/dashboard/v1/registrations/"+id, bytes.NewReader(body))
	rr := httptest.NewRecorder()

	PatchDashboardRegistration(rr, req, id)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

func TestPatchDashboardRegistration_BadBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/dashboard/v1/registrations/abc123", strings.NewReader("bad-json"))
	rr := httptest.NewRecorder()

	PatchDashboardRegistration(rr, req, "abc123")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestDeleteDashboardRegistration_Success(t *testing.T) {
	id := createTestDashboard(t)

	req := httptest.NewRequest(http.MethodDelete, "/dashboard/v1/registrations/"+id, nil)
	rr := httptest.NewRecorder()

	DeleteDashboardRegistration(rr, req, id)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected 204 No Content, got %d", rr.Code)
	}
}
