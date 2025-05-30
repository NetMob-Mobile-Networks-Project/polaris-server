package api

import (
	"net/http"
)

// GET /api/export/csv
// Description: Export all telemetry data as CSV
func (h *Handler) ExportCSVHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Generate and stream CSV file
}

// GET /api/export/kml
// Description: Export geo-locations and metadata as KML (for Google Earth)
func (h *Handler) ExportKMLHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Generate and return KML format
}
