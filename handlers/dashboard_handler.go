package handlers

import (
	"net/http"

	"github.com/amundfpl/Assignment-2/services"
	"github.com/amundfpl/Assignment-2/utils"
)

// NewDashboardHandler returns an HTTP handler function for
// GET requests to /dashboard/v1/dashboards/{id}.
// It uses the provided DashboardService to fetch dashboard data.
func NewDashboardHandler(svc services.DashboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enforce that the HTTP method must be GET
		if !utils.EnforceMethod(w, r, http.MethodGet) {
			return
		}

		// Extract dashboard ID from the URL path
		id, extractErr := utils.ExtractIDFromPath(r, utils.DashboardIDPathIndex)
		if extractErr != nil || id == "" {
			utils.WriteErrorResponse(w, utils.ErrMsgMissingOrInvalidDashboardID, http.StatusBadRequest)
			return
		}

		// Fetch populated dashboard data from the service
		dashboard, fetchErr := svc.GetPopulatedDashboardByID(id)
		if fetchErr != nil {
			utils.WriteErrorResponse(w, utils.ErrMsgDashboardFetchFailed+fetchErr.Error(), http.StatusInternalServerError)
			return
		}

		// Send the dashboard data as a successful JSON response
		utils.WriteSuccessResponse(w, dashboard, http.StatusOK)
	}
}
