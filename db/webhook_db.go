package db

import (
	"context"
	"github.com/amundfpl/Assignment-2/utils"
)

// SaveWebhook stores a new webhook in the Firestore webhook collection.
// Returns the generated document ID or an error.
func SaveWebhook(ctx context.Context, webhook utils.Webhook) (string, error) {
	docRef, _, saveErr := firestoreClient.Collection(utils.WebhookCollection).Add(ctx, webhook)
	if saveErr != nil {
		return "", saveErr
	}
	return docRef.ID, nil
}

// GetWebhookByID retrieves a single webhook from Firestore using its document ID.
// Returns the webhook or an error if not found or decoding fails.
func GetWebhookByID(ctx context.Context, id string) (*utils.Webhook, error) {
	docSnap, getErr := firestoreClient.Collection(utils.WebhookCollection).Doc(id).Get(ctx)
	if getErr != nil {
		return nil, getErr
	}

	var webhook utils.Webhook
	decodeErr := docSnap.DataTo(&webhook)
	if decodeErr != nil {
		return nil, decodeErr
	}

	webhook.ID = docSnap.Ref.ID
	return &webhook, nil
}

// GetAllWebhooks fetches all stored webhooks from Firestore.
// Returns a slice of webhook objects or an empty list if none are found.
func GetAllWebhooks(ctx context.Context) ([]utils.Webhook, error) {
	iter := firestoreClient.Collection(utils.WebhookCollection).Documents(ctx)
	var hooks []utils.Webhook

	for {
		docSnap, nextErr := iter.Next()
		if nextErr != nil {
			break // No more documents
		}
		var hook utils.Webhook
		decodeErr := docSnap.DataTo(&hook)
		if decodeErr == nil {
			hook.ID = docSnap.Ref.ID
			hooks = append(hooks, hook)
		}
	}
	return hooks, nil
}

// GetMatchingWebhooks finds webhooks registered for a specific event,
// optionally filtering by country. If a webhook has no country restriction (""), it is included.
func GetMatchingWebhooks(ctx context.Context, event, country string) ([]utils.Webhook, error) {
	iter := firestoreClient.Collection(utils.WebhookCollection).Where("event", "==", event).Documents(ctx)
	var hooks []utils.Webhook

	for {
		doc, err := iter.Next()
		if err != nil {
			break // No more documents
		}

		var hook utils.Webhook
		if err := doc.DataTo(&hook); err == nil {
			hook.ID = doc.Ref.ID
			// Match country exactly or allow wildcard (empty string)
			if hook.Country == country || hook.Country == "" {
				hooks = append(hooks, hook)
			}
		}
	}
	return hooks, nil
}

// DeleteWebhook removes a webhook from Firestore using its document ID.
// Returns an error if the operation fails.
func DeleteWebhook(ctx context.Context, id string) error {
	_, deleteErr := firestoreClient.Collection(utils.WebhookCollection).Doc(id).Delete(ctx)
	return deleteErr
}

// CountWebhooks returns the total number of webhook documents stored in Firestore.
// Returns 0 if an error occurs.
func CountWebhooks(ctx context.Context) int {
	docs, countErr := firestoreClient.Collection(utils.WebhookCollection).Documents(ctx).GetAll()
	if countErr != nil {
		return 0
	}
	return len(docs)
}
