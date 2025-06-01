package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ========== Request/Response Structs ==========

type HeatmapMetricsRequest struct {
	Lat    float64 `json:"lat" query:"lat"`       // Latitude of center point
	Lon    float64 `json:"lon" query:"lon"`       // Longitude of center point
	Radius float64 `json:"radius" query:"radius"` // Search radius in meters
}

type HeatmapPoint struct {
	Lat         float64 `json:"lat"`          // Latitude of measurement
	Lon         float64 `json:"lon"`          // Longitude of measurement
	SignalValue float64 `json:"signal_value"` // Signal strength value
	Timestamp   string  `json:"timestamp"`    // When measurement was taken
}

type TableMetricsRequest struct {
	Start time.Time `json:"start" query:"start"` // Start of time range
	End   time.Time `json:"end" query:"end"`     // End of time range
}

type TableMetricsResponse struct {
	ID          string    `json:"id"`           // Unique identifier
	Timestamp   time.Time `json:"timestamp"`    // When measurement was taken
	RSRP        float64   `json:"rsrp"`         // Reference Signal Received Power
	RSRQ        float64   `json:"rsrq"`         // Reference Signal Received Quality
	SINR        float64   `json:"sinr"`         // Signal to Interference plus Noise Ratio
	Location    string    `json:"location"`     // Human readable location
	DeviceInfo  string    `json:"device_info"`  // Device information
	NetworkInfo string    `json:"network_info"` // Network information
}

type ChartMetricsRequest struct {
	Metric string    `json:"metric" query:"metric"` // Metric to plot (rsrp, rsrq, sinr)
	Start  time.Time `json:"start" query:"start"`   // Start of time range
	End    time.Time `json:"end" query:"end"`       // End of time range
}

type ChartDataPoint struct {
	Timestamp time.Time `json:"timestamp"` // When measurement was taken
	Value     float64   `json:"value"`     // Metric value
}

type ConfigRequest struct {
	Username string `json:"username"` // Admin username
	Password string `json:"password"` // Admin password
}

type ConfigResponse struct {
	Token     string    `json:"token"`      // JWT or session token
	ExpiresAt time.Time `json:"expires_at"` // Token expiration time
}

type UpdateConfigRequest struct {
	Token string `json:"token"` // Session token to invalidate
}

// ========== Dashboard APIs ==========

// GET /api/dashboard/metrics/map?lat=...&lon=...&radius=...
// Description: Return signal quality points for heatmap visualization
// Response: List of geo-tagged signal strengths
func (h *Handler) GetHeatmapMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var req HeatmapMetricsRequest
	if err := h.parseQueryParams(r, &req); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// TODO: Query database for spatial points and signal values
	points := []HeatmapPoint{} // Replace with actual database query

	// Return response
	h.respondJSON(w, points)
}

// GET /api/dashboard/metrics/table?start=...&end=...
// Description: Return tabular list of test results in a given time range
func (h *Handler) GetTableMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var req TableMetricsRequest
	if err := h.parseQueryParams(r, &req); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// TODO: Query database for measurements in time range
	rows := []TableMetricsResponse{} // Replace with actual database query

	// Return response
	h.respondJSON(w, rows)
}

// GET /api/dashboard/metrics/chart?metric=rsrp&start=...&end=...
// Description: Return time-series values for a specific metric (for charts)
func (h *Handler) GetChartMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var req ChartMetricsRequest
	if err := h.parseQueryParams(r, &req); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	// Validate metric type
	validMetrics := map[string]bool{"rsrp": true, "rsrq": true, "sinr": true}
	if !validMetrics[req.Metric] {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	// TODO: Query database for time series data
	points := []ChartDataPoint{} // Replace with actual database query

	// Return response
	h.respondJSON(w, points)
}

// POST /api/dashboard/getconfig
// Description: Authenticate admin user and return session token or JWT
func (h *Handler) GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	var req ConfigRequest
	if err := h.parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Validate credentials
	// TODO: Generate JWT or session token

	response := ConfigResponse{
		Token:     "dummy-token", // Replace with actual token
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	h.respondJSON(w, response)
}

// POST /api/dashboard/updateconfig
// Description: Destroy user session or invalidate token
func (h *Handler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req UpdateConfigRequest
	if err := h.parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Invalidate token/session
	// TODO: Clear any server-side session data

	w.WriteHeader(http.StatusOK)
}

// Helper methods for request parsing and response handling
func (h *Handler) parseQueryParams(r *http.Request, v interface{}) error {
	query := r.URL.Query()

	switch req := v.(type) {
	case *HeatmapMetricsRequest:
		var err error
		if req.Lat, err = strconv.ParseFloat(query.Get("lat"), 64); err != nil {
			return err
		}
		if req.Lon, err = strconv.ParseFloat(query.Get("lon"), 64); err != nil {
			return err
		}
		if req.Radius, err = strconv.ParseFloat(query.Get("radius"), 64); err != nil {
			return err
		}
	case *TableMetricsRequest:
		var err error
		if req.Start, err = time.Parse(time.RFC3339, query.Get("start")); err != nil {
			return err
		}
		if req.End, err = time.Parse(time.RFC3339, query.Get("end")); err != nil {
			return err
		}
	case *ChartMetricsRequest:
		req.Metric = query.Get("metric")
		var err error
		if req.Start, err = time.Parse(time.RFC3339, query.Get("start")); err != nil {
			return err
		}
		if req.End, err = time.Parse(time.RFC3339, query.Get("end")); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported request type")
	}
	return nil
}

func (h *Handler) parseJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
