package utils

import "time"

const (
	// Routes
	DashboardRegistrationsRoute  = "/dashboard/v1/registrations/"
	DashboardRegistrationsRoute2 = "/dashboard/v1/registrations"
	DashboardDashboardsRoute     = "/dashboard/v1/dashboards/"
	DashboardNotificationsRoute  = "/dashboard/v1/notifications/"
	DashboardStatusRoute         = "/dashboard/v1/status/"
	RouteRoot                    = "/"

	// Static assets
	StaticDir       = "static"
	StaticIndexFile = "index.html"

	// Collections
	DashboardCollection     = "dashboard_configs"
	WebhookCollection       = "webhooks"
	CountryCacheCollection  = "country_cache"
	WeatherCacheCollection  = "weather_cache"
	CurrencyCacheCollection = "currency_cache"

	// Cache TTLs
	CachePurgeInterval = 1 * time.Hour
	CountryCacheTTL    = 24 * time.Hour
	WeatherCacheTTL    = 2 * time.Hour
	CurrencyCacheTTL   = 12 * time.Hour

	// Cache formatting
	WeatherCacheKeyFormat = "%.1f_%.1f"
	CacheKeySeparator     = "_"
	TimestampField        = "timestamp"
	FieldData             = "data"

	// Keys
	KeyID               = "id"
	KeyLastChange       = "lastChange"
	KeyCountry          = "country"
	KeyISOCode          = "isoCode"
	KeyFeatures         = "features"
	KeyTemperature      = "temperature"
	KeyPrecipitation    = "precipitation"
	KeyCapital          = "capital"
	KeyCoordinates      = "coordinates"
	KeyPopulation       = "population"
	KeyArea             = "area"
	KeyTargetCurrencies = "targetCurrencies"
	KeyError            = "error"

	// Config
	DefaultPort          = "8080"
	EnvPort              = "PORT"
	AddrPrefix           = ":"
	TimestampLayout      = "20060102 15:04"
	DashboardIDPathIndex = 5

	// Firebase
	FirebaseProjectID      = "ass2-cloud-refactor"
	TestFirebaseProjectID  = "ass2-cloud-refactor-test"
	CredentialsFirebaseKey = "credentials/firebase-key.json"
	CredentialsEnvVar      = "GO_FIREBASE_CREDENTIALS"
	CredentialsDir         = "credentials"
	TestCredentialsFile    = "test-serviceAccountKey.json"

	//Render environment variable
	GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS"

	// API Paths
	RESTCountriesByAlpha     = "/alpha/"
	OpenMeteoForecast        = "/v1/forecast"
	CountriesAlphaNorwayPath = "/alpha/no"
	MeteoForecastPath        = "/v1/forecast?latitude=60&longitude=10&current=temperature_2m"
	CurrencyEURToNOKPath     = "/latest?from=EUR&to=NOK"

	// API Formats
	OpenMeteoWeatherURLFmt = "%s%s?latitude=%.4f&longitude=%.4f&current=temperature_2m,precipitation"
	CurrencyAPIFmt         = "%s/%s"

	// Content Types
	ContentTypeJSON   = "application/json"
	HeaderContentType = "Content-Type"

	// Operators
	OperatorLessThan = "<"
)

// External API URLs
var (
	CurrencyAPI      = "https://api.frankfurter.app"
	RESTCountriesAPI = "https://restcountries.com/v3.1"
	OpenMeteoAPI     = "https://api.open-meteo.com"
)

// Allowed Events
var AllowedEvents = map[string]bool{
	"REGISTER": true,
	"DELETE":   true,
	"CHANGE":   true,
	"INVOKE":   true,
	"PATCH":    true,
	"LOW_TEMP": true,
}

// Webhook Events
const (
	EventLowTemp  = "LOW_TEMP"
	EventInvoke   = "INVOKE"
	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventPatch    = "PATCH"
	EventDelete   = "DELETE"
)

// Status constants
const (
	StatusVersion = "v1"
)
