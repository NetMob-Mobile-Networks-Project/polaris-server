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
	protected.Use(h.AuthMiddleware)

	// Auth APIs
	protected.HandleFunc("/auth/me", h.GetUserProfileHandler).Methods("GET")

	// Dashboard APIs
	protected.HandleFunc("/metrics/map", h.GetHeatmapMetricsHandler).Methods("GET")
	protected.HandleFunc("/metrics/table", h.GetTableMetricsHandler).Methods("GET")
	protected.HandleFunc("/metrics/chart", h.GetChartMetricsHandler).Methods("GET")

	// Export APIs
	protected.HandleFunc("/export/csv", h.ExportCSVHandler).Methods("GET")
	protected.HandleFunc("/export/kml", h.ExportKMLHandler).Methods("GET")

	// Admin routes (protected by admin middleware)
	admin := protected.PathPrefix("/admin").Subrouter()
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

func (h *Handler) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:  "error",
					Message: "Internal server error",
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware ensures the user has admin role
func (h *Handler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("user_role").(string)
		if !ok || role != "admin" {
			http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
