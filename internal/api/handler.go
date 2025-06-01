package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	database "github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage"
)

type Handler struct {
	db        *database.DB
	jwtSecret string
}

func NewHandler(db *database.DB, jwtSecret string) *Handler {
	return &Handler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Health check endpoint
	router.HandleFunc("/health", h.healthCheck).Methods("GET")

	// API version 1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Client APIs
	apiV1.HandleFunc("/client/upload", h.UploadClientDataHandler).Methods("POST")
	apiV1.HandleFunc("/client/config", h.GetClientConfigHandler).Methods("GET")

	// Dashboard APIs
	apiV1.HandleFunc("/dashboard/metrics/map", h.GetHeatmapMetricsHandler).Methods("GET")
	apiV1.HandleFunc("/dashboard/metrics/table", h.GetTableMetricsHandler).Methods("GET")
	apiV1.HandleFunc("/dashboard/metrics/chart", h.GetChartMetricsHandler).Methods("GET")
	apiV1.HandleFunc("/dashboard/getconfig", h.GetConfigHandler).Methods("GET")
	apiV1.HandleFunc("/dashboard/updateconfig", h.UpdateConfig).Methods("POST")

	// Export APIs
	apiV1.HandleFunc("/export/csv", h.ExportCSVHandler).Methods("GET")
	apiV1.HandleFunc("/export/kml", h.ExportKMLHandler).Methods("GET")

	// Auth APIs
	apiV1.HandleFunc("/auth/login", h.LoginHandler).Methods("POST")
	apiV1.HandleFunc("/auth/logout", h.LogoutHandler).Methods("POST")

	// User Management APIs
	apiV1.HandleFunc("/users", h.ListUsersHandler).Methods("GET")
	apiV1.HandleFunc("/users", h.CreateUserHandler).Methods("POST")
	apiV1.HandleFunc("/users/{id}", h.DeleteUserHandler).Methods("DELETE")
}

// healthCheck handles the health check endpoint
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
