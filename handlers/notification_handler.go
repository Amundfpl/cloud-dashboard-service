package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/services"
	"github.com/amundfpl/Assignment-2/utils"
)

// RegisterWebhook handles POST requests to register a new webhook.
// It validates input, checks event type, and stores the webhook config in the database.
func RegisterWebhook(w http.ResponseWriter, r *http.Request) {
	// Enforce correct method
	if !utils.EnforceMethod(w, r, http.MethodPost) {
		return
	}

	// Decode request body into webhook struct
	var webhook utils.Webhook
	decodeErr := json.NewDecoder(r.Body).Decode(&webhook)
	if decodeErr != nil {
		utils.WriteErrorResponse(w, utils.MsgInvalidRequestBody, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if webhook.URL == "" || webhook.Event == "" {
		utils.WriteErrorResponse(w, utils.MsgMissingWebhookFields, http.StatusBadRequest)
		return
	}

	// Normalize event and country to uppercase
	webhook.Event = strings.ToUpper(webhook.Event)
	webhook.Country = strings.ToUpper(webhook.Country)

	// Check if event type is allowed
	if !utils.AllowedEvents[webhook.Event] {
		utils.WriteErrorResponse(w, utils.MsgUnsupportedEventType+webhook.Event, http.StatusBadRequest)
		return
	}

	// Save webhook in the database
	id, saveErr := db.SaveWebhook(r.Context(), webhook)
	if saveErr != nil {
		utils.WriteErrorResponse(w, utils.MsgWebhookSaveFail, http.StatusInternalServerError)
		return
	}

	// Respond with the generated ID of the webhook
	utils.WriteIDResponse(w, id, http.StatusOK)
}

// HandleDeleteWebhook deletes a webhook by its ID.
// It ensures the ID is not empty and attempts to remove the webhook from storage.
func HandleDeleteWebhook(w http.ResponseWriter, r *http.Request, id string) {
	// Validate that ID is provided
	if id == "" {
		utils.WriteErrorResponse(w, utils.MsgMissingWebhookID, http.StatusBadRequest)
		return
	}

	// Attempt deletion
	deleteErr := services.DeleteWebhook(r.Context(), id)
	if deleteErr != nil {
		utils.WriteErrorResponse(w, utils.MsgWebhookDeleteFail+deleteErr.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}

// GetAllWebhooks returns all registered webhooks in JSON format.
// Calls the database layer to retrieve all stored webhook entries.
func GetAllWebhooks(w http.ResponseWriter, r *http.Request) {
	// Fetch all webhook entries from Firestore
	webhooks, fetchErr := db.GetAllWebhooks(r.Context())
	if fetchErr != nil {
		utils.WriteErrorResponse(w, utils.MsgWebhookFetchFail+fetchErr.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the full list
	utils.WriteSuccessResponse(w, webhooks, http.StatusOK)
}

// GetWebhookByID fetches and returns a specific webhook by ID.
// Validates the ID and retrieves the record from the database.
func GetWebhookByID(w http.ResponseWriter, r *http.Request, id string) {
	// Validate that ID is provided
	if id == "" {
		utils.WriteErrorResponse(w, utils.MsgMissingWebhookID, http.StatusBadRequest)
		return
	}

	// Fetch webhook with the given ID
	webhook, fetchErr := db.GetWebhookByID(r.Context(), id)
	if fetchErr != nil {
		utils.WriteErrorResponse(w, utils.MsgWebhookNotFound+fetchErr.Error(), http.StatusNotFound)
		return
	}

	// Respond with the found webhook
	utils.WriteSuccessResponse(w, webhook, http.StatusOK)
}
