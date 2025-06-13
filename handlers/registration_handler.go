package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/services"
	"github.com/amundfpl/Assignment-2/utils"
)

// HandleRegisterDashboard handles POST requests to register a new dashboard configuration.
// It parses the request body and delegates the registration to the service layer.
func HandleRegisterDashboard(w http.ResponseWriter, r *http.Request) {
	if !utils.EnforceMethod(w, r, http.MethodPost) {
		return
	}

	// Read and validate the request body
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		utils.WriteErrorResponse(w, utils.MsgInvalidRequestBody, http.StatusBadRequest)
		return
	}

	// Call service to register dashboard
	response, regErr := services.RegisterDashboardConfig(body)
	if regErr != nil {
		utils.WriteErrorResponse(w, utils.MsgRegisterDashboardFail+regErr.Error(), http.StatusBadRequest)
		return
	}

	// Respond with 201 Created and dashboard ID
	utils.WriteSuccessResponse(w, response, http.StatusCreated)
}

// GetRegistrationByID handles GET requests to fetch a dashboard configuration by its ID.
func GetRegistrationByID(w http.ResponseWriter, r *http.Request, id string) {
	// Retrieve config by ID from Firestore
	config, fetchErr := db.GetDashboardConfigByID(r.Context(), id)
	if fetchErr != nil {
		utils.WriteErrorResponse(w, utils.MsgDashboardNotFound+fetchErr.Error(), http.StatusNotFound)
		return
	}

	// Respond with the dashboard config
	utils.WriteSuccessResponse(w, config, http.StatusOK)
}

// GetAllRegistrations handles GET requests to retrieve all dashboard configurations.
func GetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	// Fetch all dashboard configs from DB
	dashboards, fetchErr := db.GetAllDashboardConfigs(r.Context())
	if fetchErr != nil {
		utils.WriteErrorResponse(w, utils.MsgRetrieveConfigsFail+fetchErr.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure empty array instead of null for JSON
	if dashboards == nil {
		dashboards = []utils.DashboardConfig{}
	}

	// Return list of dashboard configs
	utils.WriteSuccessResponse(w, dashboards, http.StatusOK)
}

// UpdateDashboardRegistration handles PUT requests to completely replace a configuration.
// It validates the body and forwards the request to the service.
func UpdateDashboardRegistration(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.EnforceMethod(w, r, http.MethodPut) {
		return
	}

	// Read request body
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		utils.WriteErrorResponse(w, utils.MsgInvalidRequestBody, http.StatusBadRequest)
		return
	}

	// Validate JSON before sending to service
	var temp map[string]interface{}
	if jsonErr := json.Unmarshal(body, &temp); jsonErr != nil {
		utils.WriteErrorResponse(w, utils.MsgInvalidJSON+jsonErr.Error(), http.StatusBadRequest)
		return
	}

	// Call service to update
	result, updateErr := services.UpdateDashboardConfig(r.Context(), id, body)
	if updateErr != nil {
		utils.WriteErrorResponse(w, utils.MsgUpdateConfigFail+updateErr.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated config
	utils.WriteSuccessResponse(w, result, http.StatusOK)
}

// HeadCheckDashboard handles HEAD requests to check if a dashboard exists by ID.
// This is used for lightweight existence checks.
func HeadCheckDashboard(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.EnforceMethod(w, r, http.MethodHead) {
		return
	}

	// Check if config exists
	_, fetchErr := db.GetDashboardConfigByID(r.Context(), id)
	if fetchErr != nil {
		utils.WriteErrorResponse(w, utils.MsgDashboardNotFound, http.StatusNotFound)
		return
	}

	// Return 200 OK without body
	w.WriteHeader(http.StatusOK)
}

// PatchDashboardRegistration handles PATCH requests to partially update a dashboard configuration.
func PatchDashboardRegistration(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.EnforceMethod(w, r, http.MethodPatch) {
		return
	}

	// Decode JSON patch map
	var patch map[string]interface{}
	if decodeErr := json.NewDecoder(r.Body).Decode(&patch); decodeErr != nil {
		utils.WriteErrorResponse(w, utils.MsgInvalidJSON+decodeErr.Error(), http.StatusBadRequest)
		return
	}

	// Call service to apply patch
	result, patchErr := services.PatchDashboardConfig(r.Context(), id, patch)
	if patchErr != nil {
		utils.WriteErrorResponse(w, utils.MsgPatchConfigFail+patchErr.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with updated result
	utils.WriteSuccessResponse(w, result, http.StatusOK)
}

// DeleteDashboardRegistration handles DELETE requests to remove a dashboard configuration.
func DeleteDashboardRegistration(w http.ResponseWriter, r *http.Request, id string) {
	if !utils.EnforceMethod(w, r, http.MethodDelete) {
		return
	}

	// Call service to delete by ID
	deleteErr := services.DeleteRegistrationByID(r.Context(), id)
	if deleteErr != nil {
		utils.WriteErrorResponse(w, utils.MsgDeleteConfigFail+deleteErr.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
