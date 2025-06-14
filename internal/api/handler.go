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
	apiV1.Use(h.CORSMiddleware)
	apiV1.Use(h.RecoverMiddleware)

	// Public routes
	apiV1.HandleFunc("/auth/login", h.LoginHandler).Methods("POST")

	// Client APIs (no auth required)
	apiV1.HandleFunc("/measurements/upload", h.UploadNetworkMeasurementsHandler).Methods("POST")
	apiV1.HandleFunc("/tests/http/upload", h.UploadHttpTestsHandler).Methods("POST")
	apiV1.HandleFunc("/tests/ping/upload", h.UploadPingTestsHandler).Methods("POST")
	apiV1.HandleFunc("/tests/dns/upload", h.UploadDnsTestsHandler).Methods("POST")
	apiV1.HandleFunc("/tests/web/upload", h.UploadWebTestsHandler).Methods("POST")
	apiV1.HandleFunc("/tests/sms/upload", h.UploadSmsTestsHandler).Methods("POST")
	apiV1.HandleFunc("/config", h.GetClientConfigHandler).Methods("GET")

	// Protected routes (require authentication)
	protected := apiV1.PathPrefix("").Subrouter()
	protected.Use(h.CORSMiddleware)
	protected.Use(h.AuthMiddleware)

	// Auth APIs
	protected.HandleFunc("/auth/me", h.GetUserProfileHandler).Methods("GET")

	// Dashboard APIs
	protected.HandleFunc("/metrics/avg-down-speed", h.GetAvgDownSpeedHandler).Methods("GET")
	protected.HandleFunc("/metrics/avg-up-speed", h.GetAvgUpSpeedHandler).Methods("GET")
	protected.HandleFunc("/metrics/avg-latency", h.GetAvgLatencyHandler).Methods("GET")
	protected.HandleFunc("/metrics/network-availability", h.GetNetworkAvailabilityHandler).Methods("GET")
	protected.HandleFunc("/metrics/network-histogram", h.GetNetworkHistogramHandler).Methods("GET")
	protected.HandleFunc("/metrics/network-distribution", h.GetNetworkDistributionHandler).Methods("GET")
	protected.HandleFunc("/metrics/map-data", h.GetMapDataHandler).Methods("GET")
	protected.HandleFunc("/metrics/detailed-list", h.GetDetailedListMetricsHandler).Methods("GET")
	protected.HandleFunc("/metrics/region-list", h.GetRegionListHandler).Methods("GET")

	// Export APIs
	protected.HandleFunc("/export/csv", h.ExportCSVHandler).Methods("GET")
	protected.HandleFunc("/export/kml", h.ExportKMLHandler).Methods("GET")

	// Admin routes (protected by admin middleware)
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(h.CORSMiddleware)
	admin.Use(h.AdminMiddleware)

	// Admin user management routes
	admin.HandleFunc("/users", h.ListUsersHandler).Methods("GET")
	admin.HandleFunc("/users", h.CreateUserHandler).Methods("POST")
	admin.HandleFunc("/users/{id:[0-9]+}", h.DeleteUserHandler).Methods("DELETE")

	// Admin configuration management
	admin.HandleFunc("/config", h.UpdateClientConfigHandler).Methods("PUT")
}

// healthCheck handles the health check endpoint
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
