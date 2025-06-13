package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/amundfpl/Assignment-2/db"
)

type mockWebhook struct {
	ID  string
	URL string
}

func TestTriggerWebhooks(t *testing.T) {
	// Set up a mock server to receive webhook POSTs
	var received []map[string]string
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		var payload map[string]string
		_ = json.NewDecoder(r.Body).Decode(&payload)
		mu.Lock()
		received = append(received, payload)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Insert a test webhook into Firestore
	db.GetClient().Collection("webhooks").Doc("test-webhook").Set(context.Background(), map[string]interface{}{
		"id":      "test-webhook",
		"event":   "INVOKE",
		"country": "NO",
		"url":     server.URL,
	})

	TriggerWebhooks("INVOKE", "NO")

	if len(received) == 0 {
		t.Fatal("Expected webhook to be triggered but got none")
	}

	got := received[0]
	if got["event"] != "INVOKE" {
		t.Errorf("Expected event 'INVOKE', got %s", got["event"])
	}
	if got["country"] != "NO" {
		t.Errorf("Expected country 'NO', got %s", got["country"])
	}
	if got["id"] != "test-webhook" {
		t.Errorf("Expected ID 'test-webhook', got %s", got["id"])
	}
	if got["time"] == "" {
		t.Error("Expected non-empty timestamp")
	}
}

func TestDeleteWebhook(t *testing.T) {
	ctx := context.Background()

	// Insert a test webhook
	id := "delete-me"
	_, _ = db.GetClient().Collection("webhooks").Doc(id).Set(ctx, map[string]interface{}{
		"id":      id,
		"event":   "DELETE",
		"country": "NO",
		"url":     "http://example.com",
	})

	err := DeleteWebhook(ctx, id)
	if err != nil {
		t.Fatalf("Expected no error deleting webhook, got: %v", err)
	}

	// Confirm deletion
	doc, err := db.GetClient().Collection("webhooks").Doc(id).Get(ctx)
	if err == nil && doc.Exists() {
		t.Error("Expected webhook to be deleted, but it still exists")
	}
}
