// Package services handles webhook triggering and deletion logic.
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
)

// TriggerWebhooks looks up and notifies all webhooks registered for a specific event and country.
// It builds a JSON payload with the event info and sends it to each webhook URL via HTTP POST.
func TriggerWebhooks(event, country string) {
	ctx := context.Background()

	// Fetch all webhooks that match the event and country (including wildcards).
	matchingWebhooks, fetchErr := db.GetMatchingWebhooks(ctx, event, country)
	if fetchErr != nil {
		fmt.Printf(utils.ErrFetchWebhooks, fetchErr)
		return
	}

	fmt.Printf(utils.MsgFoundWebhooks, len(matchingWebhooks), event, country)

	// Loop through the matched webhooks and trigger each one
	for _, webhook := range matchingWebhooks {
		// Prepare webhook payload
		payload := map[string]string{
			"id":      webhook.ID,
			"country": country,
			"event":   event,
			"time":    utils.CurrentTimestamp(),
		}

		// Convert to JSON
		jsonBody, marshalErr := json.Marshal(payload)
		if marshalErr != nil {
			fmt.Printf(utils.ErrMarshalWebhook, webhook.ID, marshalErr)
			continue
		}

		fmt.Printf(utils.MsgSendingWebhook, webhook.URL)

		// Send HTTP POST
		resp, postErr := http.Post(webhook.URL, utils.ContentTypeJSON, bytes.NewBuffer(jsonBody))
		if postErr != nil {
			fmt.Printf(utils.ErrSendWebhook, webhook.ID, postErr)
			continue
		}
		defer utils.CloseBody(resp.Body)

		// Log response
		fmt.Printf(utils.MsgWebhookStatus, webhook.ID, resp.Status)
	}
}

// DeleteWebhook removes a webhook with the specified ID from the Firestore database.
func DeleteWebhook(ctx context.Context, id string) error {
	return db.DeleteWebhook(ctx, id)
}
