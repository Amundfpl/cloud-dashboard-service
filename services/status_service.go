package services

import (
	"context"
	"net/http"
	"time"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
)

// serviceStartTime captures the moment the service starts (used to compute uptime)
var serviceStartTime = time.Now()

// GetSystemStatus returns an aggregated system health report including:
// - Third-party API availability (REST Countries, Open-Meteo, Currency)
// - Firestore connectivity
// - Number of registered webhooks
// - Service version
// - Uptime since start
func GetSystemStatus(ctx context.Context) utils.StatusReport {
	return utils.StatusReport{
		CountriesAPI:    checkService(utils.RESTCountriesAPI + utils.CountriesAlphaNorwayPath), // Valid ISO code
		MeteoAPI:        checkService(utils.OpenMeteoAPI + utils.MeteoForecastPath),            // Valid weather test
		CurrencyAPI:     checkService(utils.CurrencyAPI + utils.CurrencyEURToNOKPath),
		NotificationDB:  checkFirestore(ctx),
		Webhooks:        db.CountWebhooks(ctx),
		Version:         utils.StatusVersion,
		UptimeInSeconds: int64(time.Since(serviceStartTime).Seconds()),
	}
}

// checkService performs a health check against an external HTTP service.
// Returns HTTP status code if successful, or 0 if the call fails.
func checkService(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return http.StatusServiceUnavailable
	}
	defer utils.CloseBody(resp.Body)

	return resp.StatusCode
}

// checkFirestore checks Firestore connectivity.
// Returns 200 on success, or 0 on failure.
func checkFirestore(ctx context.Context) int {
	if err := db.PingFirestore(ctx); err != nil {
		return http.StatusServiceUnavailable
	}
	return http.StatusOK
}
