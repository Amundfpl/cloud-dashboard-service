package utils

// --- JSON & Request Errors ---
const (
	ErrInvalidJSONFormat     = "failed to parse registration JSON: %v"
	ErrInvalidJSONBodyFormat = "failed to parse update JSON body: %v"
	MsgInvalidJSON           = "Invalid JSON body: "
	MsgInvalidRequestBody    = "Invalid request body"
)

// --- Country / ISO Code / Config Errors ---
const (
	ErrMissingCountryOrISOCode        = "either 'country' or 'isoCode' must be provided"
	ErrInvalidISOCode                 = "no country found for ISO code: %s"
	ErrCountryResponseParseFailed     = "failed to parse REST Countries API response: %v"
	ErrInvalidCountryResp             = "invalid country response: %v"
	ErrRESTCountryFetchFailed         = "failed to fetch country from REST API: %v"
	ErrFetchCountry                   = "failed to fetch country info from REST Countries API"
	ErrFetchConfig                    = "failed to fetch dashboard config"
	ErrFetchAllConfigs                = "failed to fetch all dashboard configurations"
	ErrConfigNotFound                 = "dashboard config not found: %w"
	ErrConfigNotFoundByID             = "dashboard config not found for ID: %v"
	MsgDashboardNotFound              = "Dashboard config not found"
	ErrMsgMissingOrInvalidDashboardID = "Missing or invalid dashboard ID"
	ErrMsgDashboardFetchFailed        = "Failed to retrieve populated dashboard: "
)

// --- HTTP / API Call Errors ---
const (
	ErrHTTPGetFailed   = "HTTP GET failed: %w"
	ErrHTTPGetStatus   = "HTTP GET returned status %d"
	ErrHTTPPostFailed  = "HTTP POST failed: %w"
	ErrHTTPPostStatus  = "HTTP POST returned status %d"
	ErrHTTPReadBody    = "failed to read response body: %w"
	ErrHTTPPostMarshal = "JSON marshalling failed: %w"
)

// --- Weather & Currency ---
const (
	ErrFetchWeather        = "failed to fetch weather data: %v"
	ErrInvalidWeatherResp  = "invalid weather response structure"
	ErrFetchCurrency       = "failed to fetch currency exchange rates: %v"
	ErrInvalidCurrencyResp = "invalid currency response structure"
	ErrNoBaseCurrency      = "no base currency found"
)

// --- Enrichment Errors ---
const (
	ErrEnrichCountry  = "failed to enrich country data"
	ErrEnrichWeather  = "failed to enrich weather data"
	ErrEnrichCurrency = "failed to enrich currency data"
)

// --- Cache Errors ---
const (
	ErrCacheMiss            = "cache miss for key %s: %w"
	ErrCacheDecode          = "cache decoding error for key %s: %w"
	ErrCacheExpired         = "cache expired"
	ErrCacheMissCurrency    = "currency cache miss for key %s: %w"
	ErrCacheDecodeCurrency  = "currency cache decode error for key %s: %w"
	ErrCacheExpiredCurrency = "currency cache expired"

	ErrPurgeCountryCache  = "Country cache purge error: %v"
	ErrPurgeWeatherCache  = "Weather cache purge error: %v"
	ErrPurgeCurrencyCache = "Currency cache purge error: %v"
)

// --- Firebase / Firestore ---
const (
	ErrFirebaseInit                        = "failed to initialize Firebase app: %w"
	ErrFirestoreInit                       = "failed to initialize Firestore: %w"
	ErrInitFirebaseApp                     = "failed to initialize Firebase app: %w"
	ErrInitFirestoreClient                 = "failed to initialize Firestore client: %w"
	ErrFirestoreNotInitialized             = "Firestore client is not initialized"
	ErrFirebaseAppInitializationFailed     = "Failed to create Firebase app: %v"
	ErrFirestoreClientInitializationFailed = "Failed to initialize Firestore client: %v"
	ErrFirestoreSaveFailed                 = "failed to save dashboard config: %v"
	ErrFirestoreUpdateFailed               = "failed to update dashboard config: %v"
	ErrFirestoreDeleteFailed               = "failed to delete dashboard config: %v"
)

// --- Webhooks ---
const (
	MsgMissingWebhookFields = "Missing required fields: URL or Event"
	MsgUnsupportedEventType = "Unsupported event type: "
	MsgWebhookSaveFail      = "Failed to save webhook"
	MsgWebhookDeleteFail    = "Failed to delete webhook: "
	MsgWebhookFetchFail     = "Failed to retrieve webhooks: "
	MsgMissingWebhookID     = "Missing webhook ID"
	MsgWebhookNotFound      = "Webhook not found: "
	ErrFetchWebhooks        = "Failed to fetch webhooks: %v\n"
	ErrSendWebhook          = "Webhook %s call failed: %v\n"
	ErrMarshalWebhook       = "Failed to marshal webhook payload for %s: %v\n"
)

// --- Dashboard Flow ---
const (
	MsgRegisterDashboardFail = "Failed to register dashboard: "
	MsgRetrieveConfigsFail   = "Failed to retrieve configurations: "
	MsgUpdateConfigFail      = "Failed to update config: "
	MsgPatchConfigFail       = "Failed to patch config: "
	MsgDeleteConfigFail      = "Failed to delete registration: "
	MsgDashboardSaved        = "Firestore write successful, new doc ID:"
)

// --- Cache Purge Logging ---
const (
	MsgCachePurgeStart = "Starting cache purge..."
	MsgCachePurgeDone  = "Cache purge completed. Waiting for next cycle..."
	MsgPurgeSuccess    = "Purged %d documents from %s"
)

// --- Webhook Logging ---
const (
	MsgFoundWebhooks  = "Found %d webhooks for event=%s, country=%s\n"
	MsgSendingWebhook = "Sending webhook to: %s\n"
	MsgWebhookStatus  = "Webhook %s responded with status: %s\n"
)

// --- Server & I/O Logging ---
const (
	MsgServerStart                = "Server running on port"
	ErrMsgInitDB                  = "Could not initialize database: %v"
	ErrMsgCloseFirestore          = "Error closing Firestore: %v"
	ErrMsgServerStart             = "Failed to start server: %v"
	LogFallbackCredentialUsed     = "GO_FIREBASE_CREDENTIALS not set, using fallback: %s"
	LogWriteErrorResponseFailed   = "WriteErrorResponse: failed to write response: %v"
	LogWriteSuccessResponseFailed = "WriteSuccessResponse: failed to write response: %v"
	LogWriteIDResponseFailed      = "WriteIDResponse: failed to write response: %v"
	LogCloseBodyFailed            = "CloseBody: failed to close response body: %v"
)

// --- Misc ---
const (
	ErrMissingPathID    = "missing ID in path"
	ErrMethodNotAllowed = "Method not allowed"
)

// -- Render Errors --
const (
	ErrGOOGLE_APPLICATION_CREDENTIALS_NotSet = "GOOGLE_APPLICATION_CREDENTIALS not set — falling back to default local path"
)
