package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage/db"
)

// GET /api/v1/metrics/avg-down-speed
// Description: Get average download speed for a given time range
func (h *Handler) GetAvgDownSpeedHandler(w http.ResponseWriter, r *http.Request) {
	startQuery := r.URL.Query().Get("start")

	start := time.Now()
	end := time.Now()

	if startQuery == "last-hour" {
		start = time.Now().Add(-1 * time.Hour)
	} else if startQuery == "last-day" {
		start = time.Now().Add(-24 * time.Hour)
	} else if startQuery == "last-week" {
		start = time.Now().Add(-7 * 24 * time.Hour)
	} else if startQuery == "last-month" {
		start = time.Now().Add(-30 * 24 * time.Hour)
	}

	averageSpeed, err := h.db.GetAverageDownloadSpeed(r.Context(), db.GetAverageDownloadSpeedParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get average download speed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"averageSpeed": averageSpeed.(float64),
			"unit":         "kbps",
		},
	})
}

// GET /api/v1/metrics/avg-up-speed
// Description: Get average upload speed for a given time range
func (h *Handler) GetAvgUpSpeedHandler(w http.ResponseWriter, r *http.Request) {
	startQuery := r.URL.Query().Get("start")

	start := time.Now()
	end := time.Now()

	if startQuery == "last-hour" {
		start = time.Now().Add(-1 * time.Hour)
	} else if startQuery == "last-day" {
		start = time.Now().Add(-24 * time.Hour)
	} else if startQuery == "last-week" {
		start = time.Now().Add(-7 * 24 * time.Hour)
	} else if startQuery == "last-month" {
		start = time.Now().Add(-30 * 24 * time.Hour)
	}

	averageSpeed, err := h.db.GetAverageUploadSpeed(r.Context(), db.GetAverageUploadSpeedParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get average upload speed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"averageSpeed": averageSpeed.(float64),
			"unit":         "kbps",
		},
	})
}

// GET /api/v1/metrics/avg-latency
// Description: Get average network latency for a given time range
func (h *Handler) GetAvgLatencyHandler(w http.ResponseWriter, r *http.Request) {
	startQuery := r.URL.Query().Get("start")

	start := time.Now()
	end := time.Now()

	if startQuery == "last-hour" {
		start = time.Now().Add(-1 * time.Hour)
	} else if startQuery == "last-day" {
		start = time.Now().Add(-24 * time.Hour)
	} else if startQuery == "last-week" {
		start = time.Now().Add(-7 * 24 * time.Hour)
	} else if startQuery == "last-month" {
		start = time.Now().Add(-30 * 24 * time.Hour)
	}

	averageLatency, err := h.db.GetAverageLatency(r.Context(), db.GetAverageLatencyParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get average latency")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"averageLatency": averageLatency.(float64),
			"unit":           "ms",
		},
	})
}

// GET /api/v1/metrics/network-availability
// Description: Get network availability percentage for a given time range
func (h *Handler) GetNetworkAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	startQuery := r.URL.Query().Get("start")

	start := time.Now()
	end := time.Now()

	if startQuery == "last-hour" {
		start = time.Now().Add(-1 * time.Hour)
	} else if startQuery == "last-day" {
		start = time.Now().Add(-24 * time.Hour)
	} else if startQuery == "last-week" {
		start = time.Now().Add(-7 * 24 * time.Hour)
	} else if startQuery == "last-month" {
		start = time.Now().Add(-30 * 24 * time.Hour)
	}

	availability, err := h.db.GetNetworkAvailability(r.Context(), db.GetNetworkAvailabilityParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get network availability")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"availability": availability.(float64),
			"unit":         "percent",
		},
	})
}

// GET /api/v1/metrics/network-histogram
// Description: Get histogram data for network metrics
func (h *Handler) GetNetworkHistogramHandler(w http.ResponseWriter, r *http.Request) {
	startQuery := r.URL.Query().Get("start")

	var times []time.Time
	var labels []string
	var uploadData []float64
	var downloadData []float64
	var latencyData []float64

	start := time.Now()

	if startQuery == "last-hour" {
		start = time.Now().Add(-1 * time.Hour)
		times = []time.Time{
			start.Add(10 * time.Minute),
			start.Add(20 * time.Minute),
			start.Add(30 * time.Minute),
			start.Add(40 * time.Minute),
			start.Add(50 * time.Minute),
			start.Add(60 * time.Minute),
		}

		for _, time := range times {
			labels = append(labels, time.Format("15:04"))
		}

	} else if startQuery == "last-day" {
		start = time.Now().Add(-24 * time.Hour)
		times = []time.Time{
			start.Add(4 * time.Hour),
			start.Add(8 * time.Hour),
			start.Add(12 * time.Hour),
			start.Add(16 * time.Hour),
			start.Add(20 * time.Hour),
			start.Add(24 * time.Hour),
		}

		for _, time := range times {
			labels = append(labels, time.Format("15:04"))
		}
	} else if startQuery == "last-week" {
		start = time.Now().Add(-6 * 24 * time.Hour)
		times = []time.Time{
			start.Add(24 * time.Hour),
			start.Add(2 * 24 * time.Hour),
			start.Add(3 * 24 * time.Hour),
			start.Add(4 * 24 * time.Hour),
			start.Add(5 * 24 * time.Hour),
			start.Add(6 * 24 * time.Hour),
		}

		for _, time := range times {
			labels = append(labels, time.Format("Mon"))
		}
	} else if startQuery == "last-month" {
		start = time.Now().Add(-30 * 24 * time.Hour)
		times = []time.Time{
			start.Add(5 * 24 * time.Hour),
			start.Add(10 * 24 * time.Hour),
			start.Add(15 * 24 * time.Hour),
			start.Add(20 * 24 * time.Hour),
			start.Add(25 * 24 * time.Hour),
			start.Add(30 * 24 * time.Hour),
		}

		for _, time := range times {
			labels = append(labels, time.Format("Jan 02"))
		}
	}

	for i, time := range times {
		min := start
		max := time
		if i != 0 {
			min = times[i-1]
		}

		upload, err := h.db.GetAverageUploadSpeed(r.Context(), db.GetAverageUploadSpeedParams{
			FromTimestamp: min,
			ToTimestamp:   max,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to get network metrics")
			return
		}
		download, err := h.db.GetAverageDownloadSpeed(r.Context(), db.GetAverageDownloadSpeedParams{
			FromTimestamp: min,
			ToTimestamp:   max,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to get network metrics")
			return
		}
		latency, err := h.db.GetAverageLatency(r.Context(), db.GetAverageLatencyParams{
			FromTimestamp: min,
			ToTimestamp:   max,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to get network metrics")
			return
		}

		uploadData = append(uploadData, upload.(float64))
		downloadData = append(downloadData, download.(float64))
		latencyData = append(latencyData, latency.(float64))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"labels":       labels,
			"uploadData":   uploadData,
			"downloadData": downloadData,
			"latencyData":  latencyData,
		},
	})
}

// GET /api/v1/metrics/network-distribution
// Description: Get distribution statistics for network metrics
func (h *Handler) GetNetworkDistributionHandler(w http.ResponseWriter, r *http.Request) {
	startQuery := r.URL.Query().Get("start")

	start := time.Now()
	end := time.Now()

	if startQuery == "last-hour" {
		start = time.Now().Add(-1 * time.Hour)
	} else if startQuery == "last-day" {
		start = time.Now().Add(-24 * time.Hour)
	} else if startQuery == "last-week" {
		start = time.Now().Add(-7 * 24 * time.Hour)
	} else if startQuery == "last-month" {
		start = time.Now().Add(-30 * 24 * time.Hour)
	}

	distribution, err := h.db.GetNetworkTechnologyDistribution(r.Context(), db.GetNetworkTechnologyDistributionParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get network distribution")
		return
	}

	var labels []string
	var values []float64

	for _, d := range distribution {
		labels = append(labels, d.CellularTechnology.String)
		values = append(values, d.Percentage.(float64))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"labels": labels,
			"values": values,
		},
	})
}

// GET /api/v1/metrics/map-data
// Description: Get geo-tagged network measurements for map visualization
func (h *Handler) GetMapDataHandler(w http.ResponseWriter, r *http.Request) {
	minLat := r.URL.Query().Get("min_lat")
	maxLat := r.URL.Query().Get("max_lat")
	minLong := r.URL.Query().Get("min_long")
	maxLong := r.URL.Query().Get("max_long")

	var minLatFloat, maxLatFloat, minLongFloat, maxLongFloat float64
	if minLat != "" {
		if _, err := fmt.Sscanf(minLat, "%f", &minLatFloat); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid min_lat format")
			return
		}
	} else {
		minLatFloat = -90.0
	}
	if maxLat != "" {
		if _, err := fmt.Sscanf(maxLat, "%f", &maxLatFloat); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid max_lat format")
			return
		}
	} else {
		maxLatFloat = 90.0
	}
	if minLong != "" {
		if _, err := fmt.Sscanf(minLong, "%f", &minLongFloat); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid min_long format")
			return
		}
	} else {
		minLongFloat = -180.0
	}
	if maxLong != "" {
		if _, err := fmt.Sscanf(maxLong, "%f", &maxLongFloat); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid max_long format")
			return
		}
	} else {
		maxLongFloat = 180.0
	}

	measurements, err := h.db.GetMapData(r.Context(), db.GetMapDataParams{
		FromLatitude:  sql.NullString{String: fmt.Sprintf("%f", minLatFloat), Valid: minLatFloat != 0},
		ToLatitude:    sql.NullString{String: fmt.Sprintf("%f", maxLatFloat), Valid: maxLatFloat != 0},
		FromLongitude: sql.NullString{String: fmt.Sprintf("%f", minLongFloat), Valid: minLongFloat != 0},
		ToLongitude:   sql.NullString{String: fmt.Sprintf("%f", maxLongFloat), Valid: maxLongFloat != 0},
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get map data")
		return
	}

	// Convert measurements to response format
	var points []map[string]interface{}
	for _, m := range measurements {
		point := map[string]interface{}{
			"latitude":            m.Latitude,
			"longitude":           m.Longitude,
			"timestamp":           m.Timestamp,
			"cellular_technology": m.CellularTechnology,
			"operator":            m.NetworkOperatorName,
			"signal_strength":     m.SignalStrength,
			"signal_quality":      m.SignalQuality,
		}
		if m.AdditionalParams != nil {
			point["additional_params"] = m.AdditionalParams
		}
		points = append(points, point)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"points": points,
			"bounds": map[string]interface{}{
				"min_lat":  minLatFloat,
				"max_lat":  maxLatFloat,
				"min_long": minLongFloat,
				"max_long": maxLongFloat,
			},
		},
	})
}

// GET /api/v1/metrics/detailed-list
// Description: Get detailed list of network measurements
func (h *Handler) GetDetailedListMetricsHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := 30

	metric := r.URL.Query().Get("metric")

	pageNum := 1
	if page != "" {
		if _, err := fmt.Sscanf(page, "%d", &pageNum); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid page number")
			return
		}
	}

	var allMetrics map[string]bool = map[string]bool{
		"network": true,
		"ping":    true,
		"http":    true,
		"dns":     true,
		"web":     true,
		"sms":     true,
	}

	if !allMetrics[metric] {
		writeJSONError(w, http.StatusBadRequest, "Invalid metric")
		return
	}

	var labels []string
	var values []interface{}
	var error error

	if metric == "network" {
		list, err := h.db.ListMeasurements(r.Context(), db.ListMeasurementsParams{
			Page: pageNum,
		})
	} else if metric == "ping" {
		list, err := h.db.ListPingMeasurements(r.Context(), db.ListPingMeasurementsParams{
			Page: pageNum,
		})
	} else if metric == "http" {
		list, err := h.db.ListHttpMeasurements(r.Context(), db.ListHttpMeasurementsParams{
			Page: pageNum,
		})
	} else if metric == "dns" {
		list, err := h.db.ListDnsMeasurements(r.Context(), db.ListDnsMeasurementsParams{
			Page: pageNum,
		})
	} else if metric == "web" {
		list, err := h.db.ListWebMeasurements(r.Context(), db.ListWebMeasurementsParams{
			Page: pageNum,
		})
	} else if metric == "sms" {
		list, err := h.db.ListSmsMeasurements(r.Context(), db.ListSmsMeasurementsParams{
			Page: pageNum,
		})
	}

	if error != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get detailed list of network measurements")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data: map[string]interface{}{
			"labels": labels,
			"values": values,
		},
	})
}
