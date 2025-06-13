package handlers

import (
	"net/http"

	"github.com/amundfpl/Assignment-2/services"
	"github.com/amundfpl/Assignment-2/utils"
)

// HandleServiceStatus handles GET requests to the /dashboard/v1/status endpoint.
// It returns a JSON payload containing the system's health and uptime information.
func HandleServiceStatus(w http.ResponseWriter, r *http.Request) {
	// Enforce that the request method is GET, otherwise return 405 Method Not Allowed.
	if !utils.EnforceMethod(w, r, http.MethodGet) {
		return
	}

	// Generate a system status report including API health, DB status, uptime, etc.
	status := services.GetSystemStatus(r.Context())

	// Send the status report back to the client with HTTP 200 OK.
	utils.WriteSuccessResponse(w, status, http.StatusOK)
}
