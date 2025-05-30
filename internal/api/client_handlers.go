package api

import (
	"net/http"
)

// POST /api/data/upload
// Description: Client uploads telemetry and test results
// Body: JSON with timestamp, location, signal strength, network type, test results
// Auth: Optional token
func (h *Handler) UploadClientDataHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse JSON and store in database
}

// GET /api/config
// Description: Client fetches test configuration (intervals, test types, thresholds)
// Response: JSON with sampling interval, test list, and thresholds
func (h *Handler) GetClientConfigHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Fetch and return configuration settings from DB
}
