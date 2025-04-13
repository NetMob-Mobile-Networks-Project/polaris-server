package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler represents the API handler
type Handler struct {
	db *sql.DB
}

// NewHandler creates a new API handler
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Health check endpoint
	router.HandleFunc("/health", h.healthCheck).Methods("GET")

	// API version 1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Example endpoints
	apiV1.HandleFunc("/users", h.getUsers).Methods("GET")
	apiV1.HandleFunc("/users", h.createUser).Methods("POST")
	apiV1.HandleFunc("/users/{id}", h.getUser).Methods("GET")
	apiV1.HandleFunc("/users/{id}", h.updateUser).Methods("PUT")
	apiV1.HandleFunc("/users/{id}", h.deleteUser).Methods("DELETE")
}

// healthCheck handles the health check endpoint
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// Example handler methods (to be implemented)
func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}
