package utils

// RegistrationRequest represents the payload for creating a new dashboard registration.
type RegistrationRequest struct {
	Country  string        `json:"country"`
	ISOCode  string        `json:"isoCode"`
	Features FeatureConfig `json:"features"`
}

// DashboardConfig represents the saved configuration for a dashboard.
type DashboardConfig struct {
	ID         string        `json:"id"`
	Country    string        `json:"country"`
	ISOCode    string        `json:"isoCode"`
	Features   FeatureConfig `json:"features"`
	LastChange string        `json:"lastChange"` // Timestamp string representing last update
}

// FeatureConfig represents the optional features that can be enabled in a dashboard.
type FeatureConfig struct {
	Temperature      bool     `json:"temperature"`
	Precipitation    bool     `json:"precipitation"`
	Capital          bool     `json:"capital"`
	Coordinates      bool     `json:"coordinates"`
	Population       bool     `json:"population"`
	Area             bool     `json:"area"`
	TargetCurrencies []string `json:"targetCurrencies"` // Currency codes to compare against
}

// RegistrationResponse represents the response returned after successfully registering a dashboard.
type RegistrationResponse struct {
	ID         string `json:"id"`
	LastChange string `json:"lastChange"`
}

// DashboardResponse represents an enriched dashboard, with optional country, weather, and currency info.
type DashboardResponse struct {
	Country       string             `json:"country"`
	ISOCode       string             `json:"isoCode"`
	Capital       string             `json:"capital,omitempty"`
	Latitude      float64            `json:"latitude,omitempty"`
	Longitude     float64            `json:"longitude,omitempty"`
	Population    int                `json:"population,omitempty"`
	Area          float64            `json:"area,omitempty"`
	Temperature   float64            `json:"temperature,omitempty"`
	Precipitation float64            `json:"precipitation,omitempty"`
	ExchangeRates map[string]float64 `json:"exchangeRates,omitempty"`
}

// CountryDetails is an internal model used to represent basic country information.
type CountryDetails struct {
	Capital    string
	Latitude   float64
	Longitude  float64
	Population int
	Area       float64
}

// WeatherData contains simplified weather information retrieved from external APIs.
type WeatherData struct {
	Temperature   float64
	Precipitation float64
}

// CurrencyDetails represents currency name and symbol for a given currency code.
type CurrencyDetails struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// CountryInfoResponse represents the structure of a response from the REST Countries API.
type CountryInfoResponse struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`

	Capital    []string  `json:"capital"`
	Latlng     []float64 `json:"latlng"`
	Population int       `json:"population"`
	Area       float64   `json:"area"`
	Borders    []string  `json:"borders,omitempty"`
	Flags      struct {
		Png string `json:"png"`
	} `json:"flags"`
	Languages  map[string]string          `json:"languages"`
	Currencies map[string]CurrencyDetails `json:"currencies"`
}

// StatusResponse is the top-level structure for reporting service health in a list.
type StatusResponse struct {
	Services []ServiceStatus `json:"services"`
}

// ServiceStatus represents the status of an individual external service or dependency.
type ServiceStatus struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Status   string `json:"status"`   // "OK" or "FAIL"
	Latency  int64  `json:"latency"`  // Latency in milliseconds
	HTTPCode int    `json:"httpCode"` // HTTP status code from last check
}

// StatusReport is a compact struct used for internal monitoring and dashboard status reporting.
type StatusReport struct {
	CountriesAPI    int    `json:"countries_api"`   // HTTP status code of Countries API
	MeteoAPI        int    `json:"meteo_api"`       // HTTP status code of Open-Meteo API
	CurrencyAPI     int    `json:"currency_api"`    // HTTP status code of Currency API
	NotificationDB  int    `json:"notification_db"` // Firestore connectivity check (200 or 0)
	Webhooks        int    `json:"webhooks"`        // Number of registered webhooks
	Version         string `json:"version"`         // API version
	UptimeInSeconds int64  `json:"uptime"`          // Time since server started
}

// Notification represents a generic notification message sent to the user or client.
type Notification struct {
	Message string `json:"message"`
}

// Coordinates stores a pair of latitude and longitude values.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// PopulatedFeatures holds optional dashboard feature values populated from external services.
type PopulatedFeatures struct {
	Temperature      float64            `json:"temperature,omitempty"`
	Precipitation    float64            `json:"precipitation"`
	Capital          string             `json:"capital,omitempty"`
	Coordinates      *Coordinates       `json:"coordinates,omitempty"`
	Population       int                `json:"population,omitempty"`
	Area             float64            `json:"area,omitempty"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"`
}

// PopulatedDashboardResponse represents the full dashboard data returned by /dashboards endpoints.
type PopulatedDashboardResponse struct {
	Country       string            `json:"country"`
	ISOCode       string            `json:"isoCode"`
	Features      PopulatedFeatures `json:"features"`
	LastRetrieval string            `json:"lastRetrieval"` // Timestamp of when the data was last fetched
}

// Webhook represents a registered webhook listener.
type Webhook struct {
	ID      string `json:"id" firestore:"-"`            // Local ID, not stored in Firestore
	URL     string `json:"url" firestore:"url"`         // Target URL to POST to
	Event   string `json:"event" firestore:"event"`     // Event type: REGISTER, CHANGE, DELETE, etc.
	Country string `json:"country" firestore:"country"` // ISO code or empty for global
}
