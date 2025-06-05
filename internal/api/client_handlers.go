package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	SamplingIntervalSeconds int      `json:"sampling_interval_seconds"`
	TestTypes               []string `json:"test_types"`
}

// Request structure for updating client configuration
type UpdateConfigRequest struct {
	SamplingIntervalSeconds int      `json:"sampling_interval_seconds" validate:"required,min=60,max=3600"`
	TestTypes               []string `json:"test_types" validate:"required,min=1,dive,oneof=http ping dns web sms"`
}

// Network measurement structures
type NetworkDataUploadRequest struct {
	DeviceID     string               `json:"device_id"`
	Measurements []NetworkMeasurement `json:"measurements"`
}

type NetworkMeasurement struct {
	ID                    *int64            `json:"id,omitempty"`
	Timestamp             int64             `json:"timestamp"`
	Latitude              *float64          `json:"latitude,omitempty"`
	Longitude             *float64          `json:"longitude,omitempty"`
	CellularTechnology    *string           `json:"cellular_technology,omitempty"`
	NetworkOperatorName   *string           `json:"network_operator_name,omitempty"`
	NetworkOperatorMCCMNC *string           `json:"network_operator_mccmnc,omitempty"`
	LAC                   *int32            `json:"lac,omitempty"`
	RAC                   *int32            `json:"rac,omitempty"`
	TAC                   *int32            `json:"tac,omitempty"`
	CellID                *int64            `json:"cell_id,omitempty"`
	ARFCN                 *int32            `json:"arfcn,omitempty"`
	SignalStrength        *int32            `json:"signal_strength,omitempty"`
	SignalQuality         *int32            `json:"signal_quality,omitempty"`
	AdditionalParams      map[string]string `json:"additional_params,omitempty"`
}

// Test result structures
type TestResultUploadRequest struct {
	DeviceID string      `json:"deviceId"`
	TestType string      `json:"testType"`
	Results  interface{} `json:"results"` // Will be one of the test result types
}

type HttpTestResult struct {
	ID               *int64   `json:"id,omitempty"`
	Timestamp        int64    `json:"timestamp"`
	Success          bool     `json:"success"`
	ErrorMessage     *string  `json:"error_message,omitempty"`
	Type             string   `json:"type"` // "DOWNLOAD" or "UPLOAD"
	URL              string   `json:"url"`
	BytesTransferred *int64   `json:"bytes_transferred,omitempty"`
	DurationMillis   *int64   `json:"duration_millis,omitempty"`
	ThroughputKbps   *float64 `json:"throughput_kbps,omitempty"`
}

type PingTestResult struct {
	ID                 *int64   `json:"id,omitempty"`
	Timestamp          int64    `json:"timestamp"`
	Success            bool     `json:"success"`
	ErrorMessage       *string  `json:"error_message,omitempty"`
	Host               string   `json:"host"`
	PacketsTransmitted *int32   `json:"packets_transmitted,omitempty"`
	PacketsReceived    *int32   `json:"packets_received,omitempty"`
	MinRttMillis       *float64 `json:"min_rtt_millis,omitempty"`
	AvgRttMillis       *float64 `json:"avg_rtt_millis,omitempty"`
	MaxRttMillis       *float64 `json:"max_rtt_millis,omitempty"`
	IPAddress          *string  `json:"ip_address,omitempty"`
}

type DnsTestResult struct {
	ID                   *int64   `json:"id,omitempty"`
	Timestamp            int64    `json:"timestamp"`
	Success              bool     `json:"success"`
	ErrorMessage         *string  `json:"error_message,omitempty"`
	Hostname             string   `json:"hostname"`
	ResolutionTimeMillis *int64   `json:"resolution_time_millis,omitempty"`
	ResolvedIPAddresses  []string `json:"resolved_ip_addresses,omitempty"`
}

type WebTestResult struct {
	ID                *int64  `json:"id,omitempty"`
	Timestamp         int64   `json:"timestamp"`
	Success           bool    `json:"success"`
	ErrorMessage      *string `json:"error_message,omitempty"`
	URL               string  `json:"url"`
	LoadingTimeMillis *int64  `json:"loading_time_millis,omitempty"`
	ResponseCode      *int32  `json:"response_code,omitempty"`
	PageSizeBytes     *int64  `json:"page_size_bytes,omitempty"`
}

type SmsTestResult struct {
	ID                   *int64  `json:"id,omitempty"`
	Timestamp            int64   `json:"timestamp"`
	Success              bool    `json:"success"`
	ErrorMessage         *string `json:"error_message,omitempty"`
	DestinationAddress   string  `json:"destination_address"`
	SubmissionTimeMillis *int64  `json:"submission_time_millis,omitempty"`
	DeliveryTimeMillis   *int64  `json:"delivery_time_millis,omitempty"`
	DurationMillis       *int64  `json:"duration_millis,omitempty"`
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to store default config",
			})
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
	var req UpdateConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Validate request
	if req.SamplingIntervalSeconds < 60 || req.SamplingIntervalSeconds > 3600 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Sampling interval must be between 60 and 3600 seconds",
		})
		return
	}

	if len(req.TestTypes) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "At least one test type is required",
		})
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Invalid test type: " + testType,
			})
			return
		}
	}

	// Convert to JSON for storage
	testTypesJSON, err := json.Marshal(req.TestTypes)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to process test types",
		})
		return
	}

	// Update configuration in database
	err = h.db.UpdateClientConfig(r.Context(), db.UpdateClientConfigParams{
		SamplingIntervalSeconds: int32(req.SamplingIntervalSeconds),
		TestTypes:               testTypesJSON,
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to update configuration",
		})
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Configuration updated successfully",
		Data: ClientConfigResponse{
			SamplingIntervalSeconds: req.SamplingIntervalSeconds,
			TestTypes:               req.TestTypes,
		},
	})
}

// POST /api/v1/measurements/upload
func (h *Handler) UploadNetworkMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	var req NetworkDataUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	if req.DeviceID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Device ID is required",
		})
		return
	}

	if len(req.Measurements) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "At least one measurement is required",
		})
		return
	}

	processedCount := 0
	hasErrors := false

	// Process each measurement
	for i, m := range req.Measurements {
		// Validate timestamp
		if m.Timestamp <= 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid timestamp in measurement %d", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate coordinates if provided
		if m.Latitude != nil {
			if *m.Latitude < -90 || *m.Latitude > 90 {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid latitude in measurement %d: must be between -90 and 90", i+1),
					RecordsProcessed: &processedCount,
				})
				return
			}
		}
		if m.Longitude != nil {
			if *m.Longitude < -180 || *m.Longitude > 180 {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid longitude in measurement %d: must be between -180 and 180", i+1),
					RecordsProcessed: &processedCount,
				})
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
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid cellular technology in measurement %d: %s", i+1, *m.CellularTechnology),
					RecordsProcessed: &processedCount,
				})
				return
			}
		}

		// Validate MCCMNC format if provided (should be 5-6 digits)
		if m.NetworkOperatorMCCMNC != nil {
			if len(*m.NetworkOperatorMCCMNC) < 5 || len(*m.NetworkOperatorMCCMNC) > 6 {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid MCCMNC format in measurement %d: must be 5-6 digits", i+1),
					RecordsProcessed: &processedCount,
				})
				return
			}
			if _, err := strconv.Atoi(*m.NetworkOperatorMCCMNC); err != nil {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid MCCMNC format in measurement %d: must contain only digits", i+1),
					RecordsProcessed: &processedCount,
				})
				return
			}
		}

		// Validate signal strength if provided (typically -140 to -40 dBm)
		if m.SignalStrength != nil {
			if *m.SignalStrength < -140 || *m.SignalStrength > -30 {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid signal strength in measurement %d: must be between -140 and -30 dBm", i+1),
					RecordsProcessed: &processedCount,
				})
				return
			}
		}

		// Validate signal quality if provided (typically 0-100)
		if m.SignalQuality != nil {
			if *m.SignalQuality < 0 || *m.SignalQuality > 100 {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid signal quality in measurement %d: must be between 0 and 100", i+1),
					RecordsProcessed: &processedCount,
				})
				return
			}
		}

		// Validate ARFCN if provided (typically 0-32767)
		if m.ARFCN != nil {
			if *m.ARFCN < 0 || *m.ARFCN > 32767 {
				hasErrors = true
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(UploadResponse{
					Status:           "partial_success",
					Message:          fmt.Sprintf("Invalid ARFCN in measurement %d: must be between 0 and 32767", i+1),
					RecordsProcessed: &processedCount,
				})
				return
			}
		}

		// Validate LAC, RAC, TAC if provided (should be positive)
		if m.LAC != nil && *m.LAC < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid LAC in measurement %d: must be positive", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}
		if m.RAC != nil && *m.RAC < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid RAC in measurement %d: must be positive", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}
		if m.TAC != nil && *m.TAC < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid TAC in measurement %d: must be positive", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate Cell ID if provided (should be positive)
		if m.CellID != nil && *m.CellID < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid Cell ID in measurement %d: must be positive", i+1),
				RecordsProcessed: &processedCount,
			})
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
	// Add panic recovery
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in UploadHttpTestsHandler: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:  "error",
				Message: "Internal server error",
			})
		}
	}()

	var req struct {
		DeviceID string           `json:"deviceId"`
		TestType string           `json:"testType"`
		Results  []HttpTestResult `json:"results"`
	}

	// Log the incoming request body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Error reading request body",
		})
		return
	}
	log.Printf("Received HTTP test upload request: %s", string(body))

	// Create a new reader with the body for json.Decoder
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if req.DeviceID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Device ID is required",
		})
		return
	}

	if req.TestType == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Test type is required",
		})
		return
	}

	// Validate test type
	if req.TestType != TestTypeHTTP {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid test type: must be '%s'", TestTypeHTTP),
		})
		return
	}

	if len(req.Results) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "At least one test result is required",
		})
		return
	}

	processedCount := 0
	hasErrors := false

	for i, result := range req.Results {
		// Log each result for debugging
		log.Printf("Processing HTTP test result %d: %+v", i+1, result)

		// Validate test type
		if result.Type != "DOWNLOAD" && result.Type != "UPLOAD" {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid test type in result %d: must be 'DOWNLOAD' or 'UPLOAD'", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate URL
		if result.URL == "" {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("URL is required in result %d", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate timestamp
		if result.Timestamp <= 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid timestamp in result %d", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate throughput if provided
		if result.ThroughputKbps != nil && *result.ThroughputKbps < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid throughput in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate duration if provided
		if result.DurationMillis != nil && *result.DurationMillis < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid duration in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate bytes transferred if provided
		if result.BytesTransferred != nil && *result.BytesTransferred < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid bytes transferred in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Create database parameters with nil checks
		params := db.CreateHttpTestResultParams{
			DeviceID:  req.DeviceID,
			Timestamp: time.Unix(result.Timestamp, 0),
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

		// Log the database parameters for debugging
		log.Printf("Creating HTTP test result with params: %+v", params)

		_, err := h.db.CreateHttpTestResult(r.Context(), params)
		if err != nil {
			log.Printf("Error creating HTTP test result: %v", err)
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

	log.Printf("Successfully processed %d HTTP test results", processedCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// POST /api/v1/tests/ping/upload
func (h *Handler) UploadPingTestsHandler(w http.ResponseWriter, r *http.Request) {
	// Add panic recovery
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic in UploadPingTestsHandler: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:  "error",
				Message: "Internal server error",
			})
		}
	}()

	var req struct {
		DeviceID string           `json:"deviceId"`
		TestType string           `json:"testType"`
		Results  []PingTestResult `json:"results"`
	}

	// Log the incoming request body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Error reading request body",
		})
		return
	}
	log.Printf("Received ping test upload request: %s", string(body))

	// Create a new reader with the body for json.Decoder
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	if req.DeviceID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Device ID is required",
		})
		return
	}

	if req.TestType == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "Test type is required",
		})
		return
	}

	// Validate test type
	if req.TestType != TestTypePING {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid test type: must be '%s'", TestTypePING),
		})
		return
	}

	if len(req.Results) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadResponse{
			Status:  "error",
			Message: "At least one test result is required",
		})
		return
	}

	processedCount := 0
	hasErrors := false

	for i, result := range req.Results {
		// Log each result for debugging
		log.Printf("Processing ping test result %d: %+v", i+1, result)

		// Validate host
		if result.Host == "" {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Host is required in result %d", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate timestamp
		if result.Timestamp <= 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid timestamp in result %d", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate packets if provided
		if result.PacketsTransmitted != nil && *result.PacketsTransmitted < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid packets transmitted in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}
		if result.PacketsReceived != nil && *result.PacketsReceived < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid packets received in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Validate RTT values if provided
		if result.MinRttMillis != nil && *result.MinRttMillis < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid min RTT in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}
		if result.AvgRttMillis != nil && *result.AvgRttMillis < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid avg RTT in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}
		if result.MaxRttMillis != nil && *result.MaxRttMillis < 0 {
			hasErrors = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(UploadResponse{
				Status:           "partial_success",
				Message:          fmt.Sprintf("Invalid max RTT in result %d: must be non-negative", i+1),
				RecordsProcessed: &processedCount,
			})
			return
		}

		// Create database parameters with nil checks
		params := db.CreatePingTestResultParams{
			DeviceID:  req.DeviceID,
			Timestamp: time.Unix(result.Timestamp, 0),
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

		// Log the database parameters for debugging
		log.Printf("Creating ping test result with params: %+v", params)

		_, err := h.db.CreatePingTestResult(r.Context(), params)
		if err != nil {
			log.Printf("Error creating ping test result: %v", err)
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

	log.Printf("Successfully processed %d ping test results", processedCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadResponse{
		Status:           status,
		RecordsProcessed: &processedCount,
	})
}

// POST /api/v1/tests/dns/upload
func (h *Handler) UploadDnsTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req TestResultUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	results, ok := req.Results.([]DnsTestResult)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid test results format",
		})
		return
	}

	for _, result := range results {
		resolvedIPsJSON, _ := json.Marshal(result.ResolvedIPAddresses)
		params := db.CreateDnsTestResultParams{
			DeviceID:             req.DeviceID,
			Timestamp:            time.Unix(result.Timestamp, 0),
			Success:              result.Success,
			ErrorMessage:         sql.NullString{String: *result.ErrorMessage, Valid: result.ErrorMessage != nil},
			Hostname:             result.Hostname,
			ResolutionTimeMillis: sql.NullInt64{Int64: *result.ResolutionTimeMillis, Valid: result.ResolutionTimeMillis != nil},
			ResolvedIpAddresses:  resolvedIPsJSON,
		}

		_, err := h.db.CreateDnsTestResult(r.Context(), params)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to store DNS test result",
			})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "DNS test results uploaded successfully",
	})
}

// POST /api/v1/tests/web/upload
func (h *Handler) UploadWebTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req TestResultUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	results, ok := req.Results.([]WebTestResult)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid test results format",
		})
		return
	}

	for _, result := range results {
		params := db.CreateWebTestResultParams{
			DeviceID:          req.DeviceID,
			Timestamp:         time.Unix(result.Timestamp, 0),
			Success:           result.Success,
			ErrorMessage:      sql.NullString{String: *result.ErrorMessage, Valid: result.ErrorMessage != nil},
			Url:               result.URL,
			LoadingTimeMillis: sql.NullInt64{Int64: *result.LoadingTimeMillis, Valid: result.LoadingTimeMillis != nil},
			ResponseCode:      sql.NullInt32{Int32: *result.ResponseCode, Valid: result.ResponseCode != nil},
			PageSizeBytes:     sql.NullInt64{Int64: *result.PageSizeBytes, Valid: result.PageSizeBytes != nil},
		}

		_, err := h.db.CreateWebTestResult(r.Context(), params)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to store web test result",
			})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Web test results uploaded successfully",
	})
}

// POST /api/v1/tests/sms/upload
func (h *Handler) UploadSmsTestsHandler(w http.ResponseWriter, r *http.Request) {
	var req TestResultUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	results, ok := req.Results.([]SmsTestResult)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid test results format",
		})
		return
	}

	for _, result := range results {
		params := db.CreateSmsTestResultParams{
			DeviceID:             req.DeviceID,
			Timestamp:            time.Unix(result.Timestamp, 0),
			Success:              result.Success,
			ErrorMessage:         sql.NullString{String: *result.ErrorMessage, Valid: result.ErrorMessage != nil},
			DestinationAddress:   result.DestinationAddress,
			SubmissionTimeMillis: sql.NullInt64{Int64: *result.SubmissionTimeMillis, Valid: result.SubmissionTimeMillis != nil},
			DeliveryTimeMillis:   sql.NullInt64{Int64: *result.DeliveryTimeMillis, Valid: result.DeliveryTimeMillis != nil},
			DurationMillis:       sql.NullInt64{Int64: *result.DurationMillis, Valid: result.DurationMillis != nil},
		}

		_, err := h.db.CreateSmsTestResult(r.Context(), params)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to store SMS test result",
			})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "SMS test results uploaded successfully",
	})
}
