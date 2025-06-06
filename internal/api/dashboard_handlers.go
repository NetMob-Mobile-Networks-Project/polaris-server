package api

import (
	"net/http"
)

// GET /api/metrics/map?lat=...&lon=...&radius=...
// Description: Return signal quality points for heatmap visualization
// Response: List of geo-tagged signal strengths
func (h *Handler) GetHeatmapMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Query database for spatial points and signal values
}

// GET /api/metrics/table?start=...&end=...
// Description: Return tabular list of test results in a given time range
func (h *Handler) GetTableMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Return measurements as rows for table view
}

// GET /api/metrics/chart?metric=rsrp&start=...&end=...
// Description: Return time-series values for a specific metric (for charts)
func (h *Handler) GetChartMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Aggregate and return metric data points for plotting
}
