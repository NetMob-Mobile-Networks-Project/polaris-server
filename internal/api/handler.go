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

	// Public routes
	apiV1.HandleFunc("/auth/signup", h.SignupHandler).Methods("POST")
	apiV1.HandleFunc("/auth/login", h.LoginHandler).Methods("POST")
	apiV1.HandleFunc("/auth/logout", h.LogoutHandler).Methods("POST")

	// Protected routes (require authentication)
	protected := apiV1.PathPrefix("").Subrouter()
	protected.Use(h.AuthMiddleware)

	// Auth APIs
	protected.HandleFunc("/auth/me", h.GetUserProfileHandler).Methods("GET")

	// Client APIs
	protected.HandleFunc("/data/upload", h.UploadClientDataHandler).Methods("POST")
	protected.HandleFunc("/config", h.GetClientConfigHandler).Methods("GET")

	// Dashboard APIs
	protected.HandleFunc("/metrics/map", h.GetHeatmapMetricsHandler).Methods("GET")
	protected.HandleFunc("/metrics/table", h.GetTableMetricsHandler).Methods("GET")
	protected.HandleFunc("/metrics/chart", h.GetChartMetricsHandler).Methods("GET")

	// Export APIs
	protected.HandleFunc("/export/csv", h.ExportCSVHandler).Methods("GET")
	protected.HandleFunc("/export/kml", h.ExportKMLHandler).Methods("GET")

	// User Management APIs (admin only)
	// admin := protected.PathPrefix("").Subrouter()
	// admin.Use(h.AdminMiddleware)
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
