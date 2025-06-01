package api

import (
	"net/http"
)

// ========== Dashboard APIs ==========

// GET /api/dashboard/metrics/map?lat=...&lon=...&radius=...
// Description: Return signal quality points for heatmap visualization
// Response: List of geo-tagged signal strengths
func (h *Handler) GetHeatmapMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Query database for spatial points and signal values
}

// GET /api/dashboard/metrics/table?start=...&end=...
// Description: Return tabular list of test results in a given time range
func (h *Handler) GetTableMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Return measurements as rows for table view
}

// GET /api/dashboard/metrics/chart?metric=rsrp&start=...&end=...
// Description: Return time-series values for a specific metric (for charts)
func (h *Handler) GetChartMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Aggregate and return metric data points for plotting
}

// POST /api/dashboard/getconfig
// Description: Authenticate admin user and return session token or JWT
func (h *Handler) GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate credentials and create session
}

// POST /api/dashboard/updateconfig
// Description: Destroy user session or invalidate token
func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	// TODO: End session
}
