package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/amundfpl/Assignment-2/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ---- RegisterWebhook ----

func TestRegisterWebhook_Success(t *testing.T) {
	webhook := utils.Webhook{
		URL:     "http://example.com",
		Event:   "invoke",
		Country: "no",
	}
	body, _ := json.Marshal(webhook)

	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/notifications", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterWebhook(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"id"`) {
		t.Errorf("Expected response to contain 'id', got: %s", rr.Body.String())
	}
}

func TestRegisterWebhook_InvalidJSON(t *testing.T) {
	body := `invalid-json`
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/notifications", strings.NewReader(body))
	rr := httptest.NewRecorder()

	RegisterWebhook(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestRegisterWebhook_MissingFields(t *testing.T) {
	body := `{"url": "", "event": "", "country": "NO"}`
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/notifications", strings.NewReader(body))
	rr := httptest.NewRecorder()

	RegisterWebhook(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for missing fields, got %d", rr.Code)
	}
}

func TestRegisterWebhook_UnsupportedEvent(t *testing.T) {
	body := `{"url": "http://example.com", "event": "unknown", "country": "NO"}`
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/notifications", strings.NewReader(body))
	rr := httptest.NewRecorder()

	RegisterWebhook(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for unsupported event, got %d", rr.Code)
	}
}

func TestRegisterWebhook_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/notifications", nil)
	rr := httptest.NewRecorder()

	RegisterWebhook(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 Method Not Allowed, got %d", rr.Code)
	}
}

// ---- HandleDeleteWebhook ----

func TestHandleDeleteWebhook_Success(t *testing.T) {

	req := httptest.NewRequest(http.MethodDelete, "/dashboard/v1/notifications/test-id", nil)
	rr := httptest.NewRecorder()

	HandleDeleteWebhook(rr, req, "test-id")

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected 204 No Content, got %d", rr.Code)
	}
}

func TestHandleDeleteWebhook_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/dashboard/v1/notifications/", nil)
	rr := httptest.NewRecorder()

	HandleDeleteWebhook(rr, req, "")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for missing ID, got %d", rr.Code)
	}
}

// ---- GetAllWebhooks ----

func TestGetAllWebhooks_Success(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/notifications", nil)
	rr := httptest.NewRecorder()

	GetAllWebhooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "[") {
		t.Errorf("Expected JSON array in response, got: %s", rr.Body.String())
	}
}

// ---- GetWebhookByID ----

func TestGetWebhookByID_Success(t *testing.T) {
	// 🔨 Lag en ekte webhook først
	webhook := utils.Webhook{
		URL:     "http://example.com",
		Event:   "INVOKE",
		Country: "NO",
	}
	body, _ := json.Marshal(webhook)

	reqCreate := httptest.NewRequest(http.MethodPost, "/dashboard/v1/notifications", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := httptest.NewRecorder()

	RegisterWebhook(rrCreate, reqCreate)

	if rrCreate.Code != http.StatusOK {
		t.Fatalf("❌ Could not create webhook, got %d", rrCreate.Code)
	}

	var result map[string]string
	_ = json.Unmarshal(rrCreate.Body.Bytes(), &result)
	id := result["id"]

	// 🧪 Hent webhook med faktisk ID
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/notifications/"+id, nil)
	rr := httptest.NewRecorder()

	GetWebhookByID(rr, req, id)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"url"`) {
		t.Errorf("Expected webhook response with 'url' field, got: %s", rr.Body.String())
	}
}

func TestGetWebhookByID_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/notifications/", nil)
	rr := httptest.NewRecorder()

	GetWebhookByID(rr, req, "")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for missing ID, got %d", rr.Code)
	}
}
