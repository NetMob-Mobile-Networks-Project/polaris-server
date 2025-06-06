package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage/db"
)

// Common response structure for all handlers
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Client configuration structures
type ClientConfigResponse struct {
	SamplingIntervalSeconds int      `json:"samplingIntervalSeconds"`
	TestTypes               []string `json:"testTypes"`
}

// Upload response structure matching Kotlin client expectations
type UploadResponse struct {
	Status           string `json:"status"`            // "success", "error", "partial_success"
	Message          string `json:"message,omitempty"` // Optional error message
	RecordsProcessed *int   `json:"recordsProcessed,omitempty"`
}

// Test type constants
const (
	TestTypeHTTP = "HTTP"
	TestTypePING = "PING"
	TestTypeDNS  = "DNS"
	TestTypeWEB  = "WEB"
	TestTypeSMS  = "SMS"
)

// GET /api/v1/config
func (h *Handler) GetClientConfigHandler(w http.ResponseWriter, r *http.Request) {
	config, err := h.db.GetClientConfig(r.Context())
	if err != nil {
		defaultConfig := ClientConfigResponse{
			SamplingIntervalSeconds: 300,
			TestTypes: []string{
				"http", "ping", "dns", "web", "sms",
			},
		}
		testTypesJSON, _ := json.Marshal(defaultConfig.TestTypes)
		err = h.db.UpdateClientConfig(r.Context(), db.UpdateClientConfigParams{
			SamplingIntervalSeconds: int32(defaultConfig.SamplingIntervalSeconds),
			TestTypes:               testTypesJSON,
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to store default config")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Default configuration retrieved successfully",
			Data:    defaultConfig,
		})
		return
	}
	var testTypes []string
	json.Unmarshal([]byte(config.TestTypes), &testTypes)
	configResponse := ClientConfigResponse{
		SamplingIntervalSeconds: int(config.SamplingIntervalSeconds),
		TestTypes:               testTypes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Configuration retrieved successfully",
		Data:    configResponse,
	})
}

// PUT /api/v1/admin/config
func (h *Handler) UpdateClientConfigHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SamplingIntervalSeconds int      `json:"samplingIntervalSeconds" validate:"required,min=60,max=3600"`
		TestTypes               []string `json:"testTypes" validate:"required,min=1,dive,oneof=http ping dns web sms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.SamplingIntervalSeconds < 60 || req.SamplingIntervalSeconds > 3600 {
		writeJSONError(w, http.StatusBadRequest, "Sampling interval must be between 60 and 3600 seconds")
		return
	}

	if len(req.TestTypes) == 0 {
		writeJSONError(w, http.StatusBadRequest, "At least one test type is required")
		return
	}

	// Validate test types
	validTestTypes := map[string]bool{
		"http": true,
		"ping": true,
		"dns":  true,
		"web":  true,
		"sms":  true,
	}
	for _, testType := range req.TestTypes {
		if !validTestTypes[testType] {
			writeJSONError(w, http.StatusBadRequest, "Invalid test type: "+testType)
			return
		}
	}

	// Convert to JSON for storage
	testTypesJSON, err := json.Marshal(req.TestTypes)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to process test types")
		return
	}

	// Update configuration in database
	err = h.db.UpdateClientConfig(r.Context(), db.UpdateClientConfigParams{
		SamplingIntervalSeconds: int32(req.SamplingIntervalSeconds),
		TestTypes:               testTypesJSON,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update configuration")
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Configuration updated successfully",
		Data:    ClientConfigResponse(req),
	})
}

// POST /api/v1/measurements/upload
func (h *Handler) UploadNetworkMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID     string `json:"deviceId"`
		Measurements []struct {
			ID                    *int64            `json:"id,omitempty"`
			Timestamp             int64             `json:"timestamp"`
			Latitude              *float64          `json:"latitude,omitempty"`
			Longitude             *float64          `json:"longitude,omitempty"`
			CellularTechnology    *string           `json:"cellularTechnology,omitempty"`
			NetworkOperatorName   *string           `json:"networkOperatorName,omitempty"`
			NetworkOperatorMCCMNC *string           `json:"networkOperatorMCCMNC,omitempty"`
			LAC                   *int32            `json:"lac,omitempty"`
			RAC                   *int32            `json:"rac,omitempty"`
			TAC                   *int32            `json:"tac,omitempty"`
			CellID                *int64            `json:"cellId,omitempty"`
			ARFCN                 *int32            `json:"arfcn,omitempty"`
			SignalStrength        *int32            `json:"signalStrength,omitempty"`
			SignalQuality         *int32            `json:"signalQuality,omitempty"`
			AdditionalParams      map[string]string `json:"additionalParams,omitempty"`
		} `json:"measurements"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeUploadError(w, http.StatusBadRequest, "Invalid request body", 0)
		return
	}

	if req.DeviceID == "" {
		writeUploadError(w, http.StatusBadRequest, "Device ID is required", 0)
		return
	}

	if len(req.Measurements) == 0 {
		writeUploadError(w, http.StatusBadRequest, "At least one measurement is required", 0)
		return
	}

	processedCount := 0
	hasErrors := false

	// Process each measurement
	for i, m := range req.Measurements {
		// Validate timestamp
		if m.Timestamp <= 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid timestamp in measurement %d", i+1), processedCount)
			return
		}

		// Validate coordinates if provided
		if m.Latitude != nil {
			if *m.Latitude < -90 || *m.Latitude > 90 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid latitude in measurement %d: must be between -90 and 90", i+1), processedCount)
				return
			}
		}
		if m.Longitude != nil {
			if *m.Longitude < -180 || *m.Longitude > 180 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid longitude in measurement %d: must be between -180 and 180", i+1), processedCount)
				return
			}
		}

		// Validate cellular technology if provided
		if m.CellularTechnology != nil {
			validTechs := map[string]bool{
				"2G": true, "3G": true, "4G": true, "5G": true,
				"GSM": true, "UMTS": true, "LTE": true, "NR": true,
			}
			if !validTechs[*m.CellularTechnology] {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid cellular technology in measurement %d: %s", i+1, *m.CellularTechnology), processedCount)
				return
			}
		}

		// Validate MCCMNC format if provided
		if m.NetworkOperatorMCCMNC != nil {
			if len(*m.NetworkOperatorMCCMNC) < 5 || len(*m.NetworkOperatorMCCMNC) > 6 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid MCCMNC format in measurement %d: must be 5-6 digits", i+1), processedCount)
				return
			}
			if _, err := strconv.Atoi(*m.NetworkOperatorMCCMNC); err != nil {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid MCCMNC format in measurement %d: must contain only digits", i+1), processedCount)
				return
			}
		}

		// Validate signal strength if provided
		if m.SignalStrength != nil {
			if *m.SignalStrength < -140 || *m.SignalStrength > -30 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid signal strength in measurement %d: must be between -140 and -30 dBm", i+1), processedCount)
				return
			}
		}

		// Validate signal quality if provided
		if m.SignalQuality != nil {
			if *m.SignalQuality < -200 || *m.SignalQuality > 200 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid signal quality in measurement %d: must be between 0 and 100", i+1), processedCount)
				return
			}
		}

		// Validate ARFCN if provided
		if m.ARFCN != nil {
			if *m.ARFCN < 0 || *m.ARFCN > 32767 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid ARFCN in measurement %d: must be between 0 and 32767", i+1), processedCount)
				return
			}
		}

		// Validate LAC, RAC, TAC if provided
		if m.LAC != nil && *m.LAC < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid LAC in measurement %d: must be positive", i+1), processedCount)
			return
		}
		if m.RAC != nil && *m.RAC < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid RAC in measurement %d: must be positive", i+1), processedCount)
			return
		}
		if m.TAC != nil && *m.TAC < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid TAC in measurement %d: must be positive", i+1), processedCount)
			return
		}

		// Validate Cell ID if provided
		if m.CellID != nil && *m.CellID < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid Cell ID in measurement %d: must be positive", i+1), processedCount)
			return
		}

		additionalParamsJSON, _ := json.Marshal(m.AdditionalParams)
		params := db.CreateNetworkMeasurementParams{
			DeviceID:              req.DeviceID,
			Timestamp:             time.Unix(m.Timestamp, 0),
			Latitude:              sql.NullString{String: strconv.FormatFloat(*m.Latitude, 'f', 8, 64), Valid: m.Latitude != nil},
			Longitude:             sql.NullString{String: strconv.FormatFloat(*m.Longitude, 'f', 8, 64), Valid: m.Longitude != nil},
			CellularTechnology:    sql.NullString{String: *m.CellularTechnology, Valid: m.CellularTechnology != nil},
			NetworkOperatorName:   sql.NullString{String: *m.NetworkOperatorName, Valid: m.NetworkOperatorName != nil},
			NetworkOperatorMccmnc: sql.NullString{String: *m.NetworkOperatorMCCMNC, Valid: m.NetworkOperatorMCCMNC != nil},
			Lac:                   sql.NullInt32{Int32: *m.LAC, Valid: m.LAC != nil},
			Rac:                   sql.NullInt32{Int32: *m.RAC, Valid: m.RAC != nil},
			Tac:                   sql.NullInt32{Int32: *m.TAC, Valid: m.TAC != nil},
			CellID:                sql.NullInt64{Int64: *m.CellID, Valid: m.CellID != nil},
			Arfcn:                 sql.NullInt32{Int32: *m.ARFCN, Valid: m.ARFCN != nil},
			SignalStrength:        sql.NullInt32{Int32: *m.SignalStrength, Valid: m.SignalStrength != nil},
			SignalQuality:         sql.NullInt32{Int32: *m.SignalQuality, Valid: m.SignalQuality != nil},
			AdditionalParams:      additionalParamsJSON,
		}

		_, err := h.db.CreateNetworkMeasurement(r.Context(), params)
		if err != nil {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          "Failed to store measurement",
				RecordsProcessed: &processedCount,
			})
			return
		}

		processedCount++
	}

	// Determine final status
	status := "success"
	if hasErrors {
		status = "partial_success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// POST /api/v1/tests/http/upload
func (h *Handler) UploadHttpTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"deviceId"`
		TestType string `json:"testType"`
		Results  []struct {
			ID               *int64   `json:"id,omitempty"`
			Timestamp        int64    `json:"timestamp"`
			Success          bool     `json:"success"`
			ErrorMessage     *string  `json:"errorMessage,omitempty"`
			Type             string   `json:"type"` // "DOWNLOAD" or "UPLOAD"
			URL              string   `json:"url"`
			BytesTransferred *int64   `json:"bytesTransferred,omitempty"`
			DurationMillis   *int64   `json:"durationMillis,omitempty"`
			ThroughputKbps   *float64 `json:"throughputKbps,omitempty"`
		} `json:"results"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err), 0)
		return
	}

	if req.DeviceID == "" {
		writeUploadError(w, http.StatusBadRequest, "Device ID is required", 0)
		return
	}

	if req.TestType == "" {
		writeUploadError(w, http.StatusBadRequest, "Test type is required", 0)
		return
	}

	// Validate test type
	if req.TestType != TestTypeHTTP {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test type: must be '%s'", TestTypeHTTP), 0)
		return
	}

	if len(req.Results) == 0 {
		writeUploadError(w, http.StatusBadRequest, "At least one test result is required", 0)
		return
	}

	processedCount := 0
	hasErrors := false

	for i, result := range req.Results {
		// Validate test type
		if result.Type != "DOWNLOAD" && result.Type != "UPLOAD" {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test type in result %d: must be 'DOWNLOAD' or 'UPLOAD'", i+1), processedCount)
			return
		}

		// Validate URL
		if result.URL == "" {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("URL is required in result %d", i+1), processedCount)
			return
		}

		// Validate timestamp
		if result.Timestamp <= 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid timestamp in result %d: must be positive", i+1), processedCount)
			return
		}

		// Validate bytes transferred if provided
		if result.BytesTransferred != nil {
			if *result.BytesTransferred < 0 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid bytes transferred in result %d: must be non-negative", i+1), processedCount)
				return
			}
		}

		// Validate duration if provided
		if result.DurationMillis != nil {
			if *result.DurationMillis < 0 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid duration in result %d: must be non-negative", i+1), processedCount)
				return
			}
		}

		// Validate throughput if provided
		if result.ThroughputKbps != nil {
			if *result.ThroughputKbps < 0 {
				writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid throughput in result %d: must be non-negative", i+1), processedCount)
				return
			}
		}

		// Create database parameters with nil checks
		params := db.CreateHttpTestResultParams{
			DeviceID:  req.DeviceID,
			Timestamp: time.Unix(result.Timestamp/1000, (result.Timestamp%1000)*1000000), // Convert milliseconds to seconds and nanoseconds
			Success:   result.Success,
			TestType:  db.HttpTestResultsTestType(result.Type),
			Url:       result.URL,
		}

		// Set optional fields with nil checks
		if result.ErrorMessage != nil {
			params.ErrorMessage = sql.NullString{String: *result.ErrorMessage, Valid: true}
		}
		if result.BytesTransferred != nil {
			params.BytesTransferred = sql.NullInt64{Int64: *result.BytesTransferred, Valid: true}
		}
		if result.DurationMillis != nil {
			params.DurationMillis = sql.NullInt64{Int64: *result.DurationMillis, Valid: true}
		}
		if result.ThroughputKbps != nil {
			params.ThroughputKbps = sql.NullString{String: strconv.FormatFloat(*result.ThroughputKbps, 'f', 2, 64), Valid: true}
		}

		_, err := h.db.CreateHttpTestResult(r.Context(), params)
		if err != nil {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Failed to store HTTP test result: %v", err),
				RecordsProcessed: &processedCount,
			})
			return
		}

		processedCount++
	}

	// Determine final status
	status := "success"
	if hasErrors {
		status = "partial_success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// POST /api/v1/tests/ping/upload
func (h *Handler) UploadPingTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"deviceId"`
		TestType string `json:"testType"`
		Results  []struct {
			ID                 *int64   `json:"id,omitempty"`
			Timestamp          int64    `json:"timestamp"`
			Success            bool     `json:"success"`
			ErrorMessage       *string  `json:"errorMessage,omitempty"`
			Host               string   `json:"host"`
			PacketsTransmitted *int32   `json:"packetsTransmitted,omitempty"`
			PacketsReceived    *int32   `json:"packetsReceived,omitempty"`
			MinRttMillis       *float64 `json:"minRttMillis,omitempty"`
			AvgRttMillis       *float64 `json:"avgRttMillis,omitempty"`
			MaxRttMillis       *float64 `json:"maxRttMillis,omitempty"`
			IPAddress          *string  `json:"ipAddress,omitempty"`
		} `json:"results"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err), 0)
		return
	}

	if req.DeviceID == "" {
		writeUploadError(w, http.StatusBadRequest, "Device ID is required", 0)
		return
	}

	if req.TestType == "" {
		writeUploadError(w, http.StatusBadRequest, "Test type is required", 0)
		return
	}

	// Validate test type
	if req.TestType != TestTypePING {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test type: must be '%s'", TestTypePING), 0)
		return
	}

	if len(req.Results) == 0 {
		writeUploadError(w, http.StatusBadRequest, "At least one test result is required", 0)
		return
	}

	processedCount := 0
	hasErrors := false

	for i, result := range req.Results {
		// Validate host
		if result.Host == "" {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Host is required in result %d", i+1), processedCount)
			return
		}

		// Validate timestamp
		if result.Timestamp <= 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid timestamp in result %d: must be positive", i+1), processedCount)
			return
		}

		// Validate packets if provided
		if result.PacketsTransmitted != nil && *result.PacketsTransmitted < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid packets transmitted in result %d: must be non-negative", i+1), processedCount)
			return
		}
		if result.PacketsReceived != nil && *result.PacketsReceived < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid packets received in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Validate RTT values if provided
		if result.MinRttMillis != nil && *result.MinRttMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid min RTT in result %d: must be non-negative", i+1), processedCount)
			return
		}
		if result.AvgRttMillis != nil && *result.AvgRttMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid avg RTT in result %d: must be non-negative", i+1), processedCount)
			return
		}
		if result.MaxRttMillis != nil && *result.MaxRttMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid max RTT in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Create database parameters with nil checks
		params := db.CreatePingTestResultParams{
			DeviceID:  req.DeviceID,
			Timestamp: time.Unix(result.Timestamp/1000, (result.Timestamp%1000)*1000000), // Convert milliseconds to seconds and nanoseconds
			Success:   result.Success,
			Host:      result.Host,
		}

		// Set optional fields with nil checks
		if result.ErrorMessage != nil {
			params.ErrorMessage = sql.NullString{String: *result.ErrorMessage, Valid: true}
		}
		if result.PacketsTransmitted != nil {
			params.PacketsTransmitted = sql.NullInt32{Int32: *result.PacketsTransmitted, Valid: true}
		}
		if result.PacketsReceived != nil {
			params.PacketsReceived = sql.NullInt32{Int32: *result.PacketsReceived, Valid: true}
		}
		if result.MinRttMillis != nil {
			params.MinRttMillis = sql.NullFloat64{Float64: *result.MinRttMillis, Valid: true}
		}
		if result.AvgRttMillis != nil {
			params.AvgRttMillis = sql.NullFloat64{Float64: *result.AvgRttMillis, Valid: true}
		}
		if result.MaxRttMillis != nil {
			params.MaxRttMillis = sql.NullFloat64{Float64: *result.MaxRttMillis, Valid: true}
		}
		if result.IPAddress != nil {
			params.IpAddress = sql.NullString{String: *result.IPAddress, Valid: true}
		}

		_, err := h.db.CreatePingTestResult(r.Context(), params)
		if err != nil {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Failed to store ping test result: %v", err),
				RecordsProcessed: &processedCount,
			})
			return
		}

		processedCount++
	}

	// Determine final status
	status := "success"
	if hasErrors {
		status = "partial_success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// POST /api/v1/tests/dns/upload
func (h *Handler) UploadDnsTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"deviceId"`
		TestType string `json:"testType"`
		Results  []struct {
			ID                   *int64   `json:"id,omitempty"`
			Timestamp            int64    `json:"timestamp"`
			Success              bool     `json:"success"`
			ErrorMessage         *string  `json:"errorMessage,omitempty"`
			Hostname             string   `json:"hostname"`
			ResolutionTimeMillis *int64   `json:"resolutionTimeMillis,omitempty"`
			ResolvedIPAddresses  []string `json:"resolvedIPAddresses,omitempty"`
		} `json:"results"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err), 0)
		return
	}

	if req.DeviceID == "" {
		writeUploadError(w, http.StatusBadRequest, "Device ID is required", 0)
		return
	}

	if req.TestType == "" {
		writeUploadError(w, http.StatusBadRequest, "Test type is required", 0)
		return
	}

	// Validate test type
	if req.TestType != TestTypeDNS {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test type: must be '%s'", TestTypeDNS), 0)
		return
	}

	if len(req.Results) == 0 {
		writeUploadError(w, http.StatusBadRequest, "At least one test result is required", 0)
		return
	}

	processedCount := 0
	hasErrors := false

	for i, result := range req.Results {
		// Validate hostname
		if result.Hostname == "" {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Hostname is required in result %d", i+1), processedCount)
			return
		}

		// Validate timestamp
		if result.Timestamp <= 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid timestamp in result %d: must be positive", i+1), processedCount)
			return
		}

		// Validate resolution time if provided
		if result.ResolutionTimeMillis != nil && *result.ResolutionTimeMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid resolution time in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Create database parameters with nil checks
		params := db.CreateDnsTestResultParams{
			DeviceID:  req.DeviceID,
			Timestamp: time.Unix(result.Timestamp/1000, (result.Timestamp%1000)*1000000), // Convert milliseconds to seconds and nanoseconds
			Success:   result.Success,
			Hostname:  result.Hostname,
		}

		// Set optional fields with nil checks
		if result.ErrorMessage != nil {
			params.ErrorMessage = sql.NullString{String: *result.ErrorMessage, Valid: true}
		}
		if result.ResolutionTimeMillis != nil {
			params.ResolutionTimeMillis = sql.NullInt64{Int64: *result.ResolutionTimeMillis, Valid: true}
		}
		if result.ResolvedIPAddresses != nil {
			resolvedIPsJSON, err := json.Marshal(result.ResolvedIPAddresses)
			if err != nil {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Failed to process resolved IP addresses in result %d", i+1),
					RecordsProcessed: &processedCount,
				})
				return
			}
			params.ResolvedIpAddresses = resolvedIPsJSON
		}

		_, err := h.db.CreateDnsTestResult(r.Context(), params)
		if err != nil {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Failed to store DNS test result: %v", err),
				RecordsProcessed: &processedCount,
			})
			return
		}

		processedCount++
	}

	// Determine final status
	status := "success"
	if hasErrors {
		status = "partial_success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// POST /api/v1/tests/web/upload
func (h *Handler) UploadWebTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"deviceId"`
		TestType string `json:"testType"`
		Results  []struct {
			ID                *int64  `json:"id,omitempty"`
			Timestamp         int64   `json:"timestamp"`
			Success           bool    `json:"success"`
			ErrorMessage      *string `json:"errorMessage,omitempty"`
			URL               string  `json:"url"`
			LoadingTimeMillis *int64  `json:"loadingTimeMillis,omitempty"`
			ResponseCode      *int32  `json:"responseCode,omitempty"`
			PageSizeBytes     *int64  `json:"pageSizeBytes,omitempty"`
		} `json:"results"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err), 0)
		return
	}

	if req.DeviceID == "" {
		writeUploadError(w, http.StatusBadRequest, "Device ID is required", 0)
		return
	}

	if req.TestType == "" {
		writeUploadError(w, http.StatusBadRequest, "Test type is required", 0)
		return
	}

	// Validate test type
	if req.TestType != TestTypeWEB {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test type: must be '%s'", TestTypeWEB), 0)
		return
	}

	if len(req.Results) == 0 {
		writeUploadError(w, http.StatusBadRequest, "At least one test result is required", 0)
		return
	}

	processedCount := 0
	hasErrors := false

	for i, result := range req.Results {
		// Validate URL
		if result.URL == "" {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("URL is required in result %d", i+1), processedCount)
			return
		}

		// Validate timestamp
		if result.Timestamp <= 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid timestamp in result %d: must be positive", i+1), processedCount)
			return
		}

		// Validate loading time if provided
		if result.LoadingTimeMillis != nil && *result.LoadingTimeMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid loading time in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Validate response code if provided
		if result.ResponseCode != nil && (*result.ResponseCode < 100 || *result.ResponseCode > 599) {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid response code in result %d: must be between 100 and 599", i+1), processedCount)
			return
		}

		// Validate page size if provided
		if result.PageSizeBytes != nil && *result.PageSizeBytes < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid page size in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Create database parameters with nil checks
		params := db.CreateWebTestResultParams{
			DeviceID:  req.DeviceID,
			Timestamp: time.Unix(result.Timestamp/1000, (result.Timestamp%1000)*1000000), // Convert milliseconds to seconds and nanoseconds
			Success:   result.Success,
			Url:       result.URL,
		}

		// Set optional fields with nil checks
		if result.ErrorMessage != nil {
			params.ErrorMessage = sql.NullString{String: *result.ErrorMessage, Valid: true}
		}
		if result.LoadingTimeMillis != nil {
			params.LoadingTimeMillis = sql.NullInt64{Int64: *result.LoadingTimeMillis, Valid: true}
		}
		if result.ResponseCode != nil {
			params.ResponseCode = sql.NullInt32{Int32: *result.ResponseCode, Valid: true}
		}
		if result.PageSizeBytes != nil {
			params.PageSizeBytes = sql.NullInt64{Int64: *result.PageSizeBytes, Valid: true}
		}

		_, err := h.db.CreateWebTestResult(r.Context(), params)
		if err != nil {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Failed to store web test result: %v", err),
				RecordsProcessed: &processedCount,
			})
			return
		}

		processedCount++
	}

	// Determine final status
	status := "success"
	if hasErrors {
		status = "partial_success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// POST /api/v1/tests/sms/upload
func (h *Handler) UploadSmsTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"deviceId"`
		TestType string `json:"testType"`
		Results  []struct {
			ID                   *int64  `json:"id,omitempty"`
			Timestamp            int64   `json:"timestamp"`
			Success              bool    `json:"success"`
			ErrorMessage         *string `json:"errorMessage,omitempty"`
			DestinationAddress   string  `json:"destinationAddress"`
			SubmissionTimeMillis *int64  `json:"submissionTimeMillis,omitempty"`
			DeliveryTimeMillis   *int64  `json:"deliveryTimeMillis,omitempty"`
			DurationMillis       *int64  `json:"durationMillis,omitempty"`
		} `json:"results"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err), 0)
		return
	}

	if req.DeviceID == "" {
		writeUploadError(w, http.StatusBadRequest, "Device ID is required", 0)
		return
	}

	if req.TestType == "" {
		writeUploadError(w, http.StatusBadRequest, "Test type is required", 0)
		return
	}

	// Validate test type
	if req.TestType != TestTypeSMS {
		writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test type: must be '%s'", TestTypeSMS), 0)
		return
	}

	if len(req.Results) == 0 {
		writeUploadError(w, http.StatusBadRequest, "At least one test result is required", 0)
		return
	}

	processedCount := 0
	hasErrors := false

	for i, result := range req.Results {
		// Validate destination address
		if result.DestinationAddress == "" {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Destination address is required in result %d", i+1), processedCount)
			return
		}

		// Validate timestamp
		if result.Timestamp <= 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid timestamp in result %d: must be positive", i+1), processedCount)
			return
		}

		// Validate submission time if provided
		if result.SubmissionTimeMillis != nil && *result.SubmissionTimeMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid submission time in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Validate delivery time if provided
		if result.DeliveryTimeMillis != nil && *result.DeliveryTimeMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid delivery time in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Validate duration if provided
		if result.DurationMillis != nil && *result.DurationMillis < 0 {
			writeUploadError(w, http.StatusBadRequest, fmt.Sprintf("Invalid duration in result %d: must be non-negative", i+1), processedCount)
			return
		}

		// Create database parameters with nil checks
		params := db.CreateSmsTestResultParams{
			DeviceID:           req.DeviceID,
			Timestamp:          time.Unix(result.Timestamp/1000, (result.Timestamp%1000)*1000000), // Convert milliseconds to seconds and nanoseconds
			Success:            result.Success,
			DestinationAddress: result.DestinationAddress,
		}

		// Set optional fields with nil checks
		if result.ErrorMessage != nil {
			params.ErrorMessage = sql.NullString{String: *result.ErrorMessage, Valid: true}
		}
		if result.SubmissionTimeMillis != nil {
			params.SubmissionTimeMillis = sql.NullInt64{Int64: *result.SubmissionTimeMillis, Valid: true}
		}
		if result.DeliveryTimeMillis != nil {
			params.DeliveryTimeMillis = sql.NullInt64{Int64: *result.DeliveryTimeMillis, Valid: true}
		}
		if result.DurationMillis != nil {
			params.DurationMillis = sql.NullInt64{Int64: *result.DurationMillis, Valid: true}
		}

		_, err := h.db.CreateSmsTestResult(r.Context(), params)
		if err != nil {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Failed to store SMS test result: %v", err),
				RecordsProcessed: &processedCount,
			})
			return
		}

		processedCount++
	}

	// Determine final status
	status := "success"
	if hasErrors {
		status = "partial_success"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// Helper to write a JSON error response
func writeJSONError(w http.ResponseWriter, status int, message string) {
	log.Printf("Upload error: %s", message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: message,
	})
}

// Helper to write an upload error response
func writeUploadError(w http.ResponseWriter, status int, message string, processedCount int) {
	log.Printf("Upload error: %s", message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           "partial_success",
		Message:          message,
		RecordsProcessed: &processedCount,
	})
}
