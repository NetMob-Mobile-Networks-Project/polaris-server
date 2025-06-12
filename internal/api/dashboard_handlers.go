package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get average download speed: %v", err))
		return
	}

	averageSpeedFloat, _ := strconv.ParseFloat(averageSpeed, 64)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Average download speed",
		Data: map[string]interface{}{
			"averageSpeed": averageSpeedFloat,
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
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get average upload speed: %v", err))
		return
	}

	averageSpeedFloat, _ := strconv.ParseFloat(averageSpeed, 64)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Average upload speed",
		Data: map[string]interface{}{
			"averageSpeed": averageSpeedFloat,
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
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get average latency: %v", err))
		return
	}

	averageLatencyFloat, _ := strconv.ParseFloat(averageLatency, 64)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Average network latency",
		Data: map[string]interface{}{
			"averageLatency": averageLatencyFloat,
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

	httpResult, err := h.db.GetHttpSuccessAndTotal(r.Context(), db.GetHttpSuccessAndTotalParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get http success and total: %v", err))
		return
	}

	pingResult, err := h.db.GetPingSuccessAndTotal(r.Context(), db.GetPingSuccessAndTotalParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get ping success and total: %v", err))
		return
	}

	dnsResult, err := h.db.GetDnsSuccessAndTotal(r.Context(), db.GetDnsSuccessAndTotalParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get dns success and total: %v", err))
		return
	}

	webResult, err := h.db.GetWebSuccessAndTotal(r.Context(), db.GetWebSuccessAndTotalParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get web success and total: %v", err))
		return
	}

	smsResult, err := h.db.GetSmsSuccessAndTotal(r.Context(), db.GetSmsSuccessAndTotalParams{
		FromTimestamp: start,
		ToTimestamp:   end,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get sms success and total: %v", err))
		return
	}

	availability := float64(httpResult.HttpSuccess+
		pingResult.PingSuccess+
		dnsResult.DnsSuccess+
		webResult.WebSuccess+
		smsResult.SmsSuccess+1) / float64(httpResult.HttpTotal+
		pingResult.PingTotal+
		dnsResult.DnsTotal+
		webResult.WebTotal+
		smsResult.SmsTotal+1)

	availabilityFloat := availability * 100

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Network availability percentage",
		Data: map[string]interface{}{
			"availability": availabilityFloat,
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

		for _, _time := range times {
			labels = append(labels, _time.Add(3*time.Hour+30*time.Minute).Format("15:04"))
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

		for _, _time := range times {
			labels = append(labels, _time.Add(3*time.Hour+30*time.Minute).Format("15:04"))
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

		for _, _time := range times {
			labels = append(labels, _time.Add(3*time.Hour+30*time.Minute).Format("Mon"))
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

		for _, _time := range times {
			labels = append(labels, _time.Add(3*time.Hour+30*time.Minute).Format("Jan 02"))
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
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get network metrics: %v", err))
			return
		}
		download, err := h.db.GetAverageDownloadSpeed(r.Context(), db.GetAverageDownloadSpeedParams{
			FromTimestamp: min,
			ToTimestamp:   max,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get network metrics: %v", err))
			return
		}
		latency, err := h.db.GetAverageLatency(r.Context(), db.GetAverageLatencyParams{
			FromTimestamp: min,
			ToTimestamp:   max,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get network metrics: %v", err))
			return
		}

		uploadFloat, _ := strconv.ParseFloat(upload, 64)
		downloadFloat, _ := strconv.ParseFloat(download, 64)
		latencyFloat, _ := strconv.ParseFloat(latency, 64)

		uploadData = append(uploadData, uploadFloat)
		downloadData = append(downloadData, downloadFloat)
		latencyData = append(latencyData, latencyFloat)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Network metrics histogram data",
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
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get network distribution: %v", err))
		return
	}

	sum := 0.0

	for _, v := range distribution {
		sum += float64(v.CellularTechnologyCount)
	}

	if sum == 0 {
		sum += 0.1
	}

	var labels []string
	var values []float64

	for _, d := range distribution {
		labels = append(labels, d.CellularTechnology.String)
		values = append(values, float64(d.CellularTechnologyCount)/sum)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Network technology distribution",
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
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get map data: %v", err))
		return
	}

	// Convert measurements to response format
	var points []map[string]interface{}
	for _, m := range measurements {
		point := map[string]interface{}{
			"id":                      m.ID,
			"device_id":               m.DeviceID,
			"timestamp":               m.Timestamp,
			"latitude":                NullStringToString(m.Latitude),
			"longitude":               NullStringToString(m.Longitude),
			"cellular_technology":     NullStringToString(m.CellularTechnology),
			"network_operator_name":   NullStringToString(m.NetworkOperatorName),
			"network_operator_mccmnc": NullStringToString(m.NetworkOperatorMccmnc),
			"lac":                     NullInt32ToInt(m.Lac),
			"rac":                     NullInt32ToInt(m.Rac),
			"tac":                     NullInt32ToInt(m.Tac),
			"cell_id":                 NullInt64ToInt(m.CellID),
			"arfcn":                   NullInt32ToInt(m.Arfcn),
			"signal_strength":         NullInt32ToInt(m.SignalStrength),
			"signal_quality":          NullInt32ToInt(m.SignalQuality),
			"additional_params":       RawMessageToString(m.AdditionalParams),
		}
		points = append(points, point)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Geo-tagged network measurements",
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

	page := r.URL.Query().Get("page")
	limit := int32(30)

	metric := r.URL.Query().Get("metric")

	pageNum := int32(1)
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
	var values []map[string]interface{}

	if metric == "network" {
		list, err := h.db.ListMeasurements(r.Context(), db.ListMeasurementsParams{
			Limit:         limit,
			Offset:        (pageNum - 1) * limit,
			FromTimestamp: start,
			ToTimestamp:   end,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get detailed list of network measurements: %v", err))
			return
		}

		for _, m := range list {
			labels = []string{
				"timestamp",
				"device_id",
				"cellular_technology",
				"network_operator_name",
				"network_operator_mccmnc",
				"lac",
				"rac",
				"tac",
				"cell_id",
				"arfcn",
				"signal_strength",
				"signal_quality",
				"latitude",
				"longitude",
				"additional_params",
			}

			values = append(values, map[string]interface{}{
				"timestamp":               m.Timestamp,
				"device_id":               m.DeviceID,
				"cellular_technology":     NullStringToString(m.CellularTechnology),
				"network_operator_name":   NullStringToString(m.NetworkOperatorName),
				"network_operator_mccmnc": NullStringToString(m.NetworkOperatorMccmnc),
				"lac":                     NullInt32ToInt(m.Lac),
				"rac":                     NullInt32ToInt(m.Rac),
				"tac":                     NullInt32ToInt(m.Tac),
				"cell_id":                 NullInt64ToInt(m.CellID),
				"arfcn":                   NullInt32ToInt(m.Arfcn),
				"signal_strength":         NullInt32ToInt(m.SignalStrength),
				"signal_quality":          NullInt32ToInt(m.SignalQuality),
				"latitude":                NullStringToString(m.Latitude),
				"longitude":               NullStringToString(m.Longitude),
				"additional_params":       RawMessageToString(m.AdditionalParams),
			})
		}
	} else if metric == "ping" {
		list, err := h.db.ListPingMeasurements(r.Context(), db.ListPingMeasurementsParams{
			Limit:         limit,
			Offset:        (pageNum - 1) * limit,
			FromTimestamp: start,
			ToTimestamp:   end,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get detailed list of network measurements: %v", err))
			return
		}

		labels = []string{
			"timestamp",
			"device_id",
			"success",
			"error_message",
			"host",
			"packets_transmitted",
			"packets_received",
			"min_rtt_millis",
			"avg_rtt_millis",
			"max_rtt_millis",
			"ip_address",
		}

		for _, m := range list {
			values = append(values, map[string]interface{}{
				"timestamp":           m.Timestamp,
				"device_id":           m.DeviceID,
				"success":             m.Success,
				"error_message":       NullStringToString(m.ErrorMessage),
				"host":                m.Host,
				"packets_transmitted": NullInt32ToInt(m.PacketsTransmitted),
				"packets_received":    NullInt32ToInt(m.PacketsReceived),
				"min_rtt_millis":      NullFloat64ToFloat(m.MinRttMillis),
				"avg_rtt_millis":      NullFloat64ToFloat(m.AvgRttMillis),
				"max_rtt_millis":      NullFloat64ToFloat(m.MaxRttMillis),
				"ip_address":          NullStringToString(m.IpAddress),
			})
		}
	} else if metric == "http" {
		list, err := h.db.ListHttpMeasurements(r.Context(), db.ListHttpMeasurementsParams{
			Limit:         limit,
			Offset:        (pageNum - 1) * limit,
			FromTimestamp: start,
			ToTimestamp:   end,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to get detailed list of network measurements")
			return
		}

		labels = []string{
			"timestamp",
			"device_id",
			"success",
			"error_message",
			"test_type",
			"url",
			"bytes_transferred",
			"duration_millis",
			"throughput_kbps",
		}

		for _, m := range list {
			values = append(values, map[string]interface{}{
				"timestamp":         m.Timestamp,
				"device_id":         m.DeviceID,
				"success":           m.Success,
				"error_message":     NullStringToString(m.ErrorMessage),
				"test_type":         string(m.TestType),
				"url":               m.Url,
				"bytes_transferred": NullInt64ToInt(m.BytesTransferred),
				"duration_millis":   NullInt64ToInt(m.DurationMillis),
				"throughput_kbps":   NullStringToString(m.ThroughputKbps),
			})
		}
	} else if metric == "dns" {
		list, err := h.db.ListDnsMeasurements(r.Context(), db.ListDnsMeasurementsParams{
			Limit:         limit,
			Offset:        (pageNum - 1) * limit,
			FromTimestamp: start,
			ToTimestamp:   end,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get detailed list of network measurements: %v", err))
			return
		}

		labels = []string{
			"timestamp",
			"device_id",
			"success",
			"error_message",
			"hostname",
			"resolution_time_millis",
			"resolved_ip_addresses",
		}

		for _, m := range list {
			values = append(values, map[string]interface{}{
				"timestamp":              m.Timestamp,
				"device_id":              m.DeviceID,
				"success":                m.Success,
				"error_message":          NullStringToString(m.ErrorMessage),
				"hostname":               m.Hostname,
				"resolution_time_millis": NullInt64ToInt(m.ResolutionTimeMillis),
				"resolved_ip_addresses":  RawMessageToString(m.ResolvedIpAddresses),
			})
		}
	} else if metric == "web" {
		list, err := h.db.ListWebMeasurements(r.Context(), db.ListWebMeasurementsParams{
			Limit:         limit,
			Offset:        (pageNum - 1) * limit,
			FromTimestamp: start,
			ToTimestamp:   end,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to get detailed list of network measurements")
			return
		}

		labels = []string{
			"timestamp",
			"device_id",
			"success",
			"error_message",
			"url",
			"loading_time_millis",
			"response_code",
			"page_size_bytes",
		}

		for _, m := range list {
			values = append(values, map[string]interface{}{
				"timestamp":           m.Timestamp,
				"device_id":           m.DeviceID,
				"success":             m.Success,
				"error_message":       NullStringToString(m.ErrorMessage),
				"url":                 m.Url,
				"loading_time_millis": NullInt64ToInt(m.LoadingTimeMillis),
				"response_code":       NullInt32ToInt(m.ResponseCode),
				"page_size_bytes":     NullInt64ToInt(m.PageSizeBytes),
			})
		}
	} else if metric == "sms" {
		list, err := h.db.ListSmsMeasurements(r.Context(), db.ListSmsMeasurementsParams{
			Limit:         limit,
			Offset:        (pageNum - 1) * limit,
			FromTimestamp: start,
			ToTimestamp:   end,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to get detailed list of network measurements")
			return
		}

		labels = []string{
			"timestamp",
			"device_id",
			"success",
			"error_message",
			"destination_address",
			"submission_time_millis",
			"delivery_time_millis",
			"duration_millis",
			"created_at",
		}

		for _, m := range list {
			values = append(values, map[string]interface{}{
				"timestamp":              m.Timestamp,
				"device_id":              m.DeviceID,
				"success":                m.Success,
				"error_message":          NullStringToString(m.ErrorMessage),
				"destination_address":    m.DestinationAddress,
				"submission_time_millis": NullInt64ToInt(m.SubmissionTimeMillis),
				"delivery_time_millis":   NullInt64ToInt(m.DeliveryTimeMillis),
				"duration_millis":        NullInt64ToInt(m.DurationMillis),
				"created_at":             m.CreatedAt,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: fmt.Sprintf("Detailed list of %s measurements", metric),
		Data: map[string]interface{}{
			"labels": labels,
			"values": values,
		},
	})
}
