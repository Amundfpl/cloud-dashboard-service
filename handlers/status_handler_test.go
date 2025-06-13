package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleServiceStatus(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/status", nil)
	rr := httptest.NewRecorder()

	HandleServiceStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	body := rr.Body.String()

	// Sjekk minst én expected field
	if !strings.Contains(body, "currency_api") || !strings.Contains(body, "uptime") {
		t.Errorf("Expected response to contain 'currency_api' and 'uptime', got: %s", body)
	}
}
