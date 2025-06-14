package server

import (
	"github.com/amundfpl/Assignment-2/handlers"
	"github.com/amundfpl/Assignment-2/services"
	"github.com/amundfpl/Assignment-2/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// InitializeRoutes sets up all HTTP route handlers and returns a configured http.Handler.
// This includes registration, webhook, dashboard, status, and static file routes.
func InitializeRoutes() http.Handler {
	realService := services.RealDashboardService{}
	getOneHandler := handlers.NewDashboardHandler(realService)
	router := http.NewServeMux()

	// Dashboard registration endpoints
	router.HandleFunc(utils.DashboardRegistrationsRoute2, registrationsDispatcher)
	router.HandleFunc(utils.DashboardRegistrationsRoute, registrationsDispatcher)

	// Dashboard visualization endpoints
	router.HandleFunc(utils.DashboardDashboardsRoute, getOneHandler)

	// Webhook notification endpoints
	router.HandleFunc(utils.DashboardNotificationsRoute, notificationsDispatcher)

	// Status check endpoint
	router.HandleFunc(utils.DashboardStatusRoute, handlers.HandleServiceStatus)

	// Serve static content and homepage fallback
	router.HandleFunc(utils.RouteRoot, FileServerWithFallback(utils.StaticDir))
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return router
}

// registrationsDispatcher handles all registration-related methods (POST, GET, PUT, PATCH, etc.).
func registrationsDispatcher(w http.ResponseWriter, r *http.Request) {
	basePath := utils.DashboardRegistrationsRoute2
	trimmedPath := strings.TrimPrefix(r.URL.Path, basePath)
	trimmedPath = strings.Trim(trimmedPath, "/")
	id := trimmedPath

	if id == "" {
		id = r.URL.Query().Get("id")
	}

	switch {
	case id == "":
		switch r.Method {
		case http.MethodGet:
			handlers.GetAllRegistrations(w, r)
		case http.MethodPost:
			handlers.HandleRegisterDashboard(w, r)
		default:
			http.Error(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	default:
		switch r.Method {
		case http.MethodGet:
			handlers.GetRegistrationByID(w, r, id)
		case http.MethodPut:
			handlers.UpdateDashboardRegistration(w, r, id)
		case http.MethodPatch:
			handlers.PatchDashboardRegistration(w, r, id)
		case http.MethodDelete:
			handlers.DeleteDashboardRegistration(w, r, id)
		case http.MethodHead:
			handlers.HeadCheckDashboard(w, r, id)
		default:
			http.Error(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	}
}

// notificationsDispatcher handles webhook registration and deletion endpoints.
func notificationsDispatcher(w http.ResponseWriter, r *http.Request) {
	basePath := utils.DashboardNotificationsRoute
	trimmedPath := strings.TrimPrefix(r.URL.Path, basePath)
	trimmedPath = strings.Trim(trimmedPath, "/")
	id := trimmedPath

	switch {
	case id == "":
		switch r.Method {
		case http.MethodPost:
			handlers.RegisterWebhook(w, r)
		case http.MethodGet:
			handlers.GetAllWebhooks(w, r)
		default:
			http.Error(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	default:
		switch r.Method {
		case http.MethodGet:
			handlers.GetWebhookByID(w, r, id)
		case http.MethodDelete:
			handlers.HandleDeleteWebhook(w, r, id)
		default:
			http.Error(w, utils.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		}
	}
}

// FileServerWithFallback serves static files from the provided directory.
// If a file is not found, it falls back to serving "index.html".
func FileServerWithFallback(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		if requestPath == utils.RouteRoot {
			http.ServeFile(w, r, filepath.Join(dir, utils.StaticIndexFile))
			return
		}

		fullPath := filepath.Join(dir, requestPath)
		if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
			http.Redirect(w, r, utils.RouteRoot, http.StatusFound)
			return
		}

		http.ServeFile(w, r, fullPath)
	}
}
