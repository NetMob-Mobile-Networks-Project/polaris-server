package api

import (
	"database/sql"
	"encoding/json"
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
	SamplingIntervalSeconds int                       `json:"sampling_interval_seconds"`
	TestTypes               []string                  `json:"test_types"`
	Thresholds              map[string]ThresholdRange `json:"thresholds"`
}

type ThresholdRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// Measurement structures
type MeasurementRequest struct {
	Timestamp     time.Time `json:"timestamp"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	NetworkType   string    `json:"network_type"`
	PLMNID        *string   `json:"plmn_id,omitempty"`
	LAC           *string   `json:"lac,omitempty"`
	TAC           *string   `json:"tac,omitempty"`
	CellID        *string   `json:"cell_id,omitempty"`
	ARFCN         *int32    `json:"arfcn,omitempty"`
	FrequencyBand *string   `json:"frequency_band,omitempty"`
	RSSI          *int32    `json:"rssi,omitempty"`
	EcNo          *float64  `json:"ec_no,omitempty"`
	RSCP          *int32    `json:"rscp,omitempty"`
	RSRQ          *float64  `json:"rsrq,omitempty"`
	RSRP          *int32    `json:"rsrp,omitempty"`
	DownloadSpeed *float64  `json:"download_speed,omitempty"`
	UploadSpeed   *float64  `json:"upload_speed,omitempty"`
	PingMS        *int32    `json:"ping_ms,omitempty"`
	DNSMS         *int32    `json:"dns_ms,omitempty"`
	HTTPLatencyMS *int32    `json:"http_latency_ms,omitempty"`
	SMSDeliveryMS *int32    `json:"sms_delivery_ms,omitempty"`
}

type MeasurementResponse struct {
	ID int64 `json:"id"`
}

// GET /api/config
// Description: Client fetches test configuration (intervals, test types, thresholds)
// Response: JSON with success status, message, and config data
func (h *Handler) GetClientConfigHandler(w http.ResponseWriter, r *http.Request) {
	config, err := h.db.GetClientConfig(r.Context())
	if err != nil {
		defaultConfig := ClientConfigResponse{
			SamplingIntervalSeconds: 300,
			TestTypes: []string{
				"signal_strength", "download_speed", "upload_speed", "ping", "dns", "http", "sms",
			},
			Thresholds: map[string]ThresholdRange{
				"signal_strength": {Min: -90, Max: -50},
				"download_speed":  {Min: 1, Max: 100},
				"upload_speed":    {Min: 1, Max: 50},
				"ping":            {Min: 0, Max: 100},
				"dns":             {Min: 0, Max: 200},
				"http":            {Min: 0, Max: 500},
				"sms":             {Min: 0, Max: 1000},
			},
		}
		testTypesJSON, _ := json.Marshal(defaultConfig.TestTypes)
		thresholdsJSON, _ := json.Marshal(defaultConfig.Thresholds)
		err = h.db.UpdateClientConfig(r.Context(), db.UpdateClientConfigParams{
			SamplingIntervalSeconds: int32(defaultConfig.SamplingIntervalSeconds),
			TestTypes:               testTypesJSON,
			Thresholds:              thresholdsJSON,
		})
		if err != nil {
			http.Error(w, "Failed to store default config", http.StatusInternalServerError)
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
	var thresholds map[string]ThresholdRange
	json.Unmarshal([]byte(config.TestTypes), &testTypes)
	json.Unmarshal([]byte(config.Thresholds), &thresholds)
	configResponse := ClientConfigResponse{
		SamplingIntervalSeconds: int(config.SamplingIntervalSeconds),
		TestTypes:               testTypes,
		Thresholds:              thresholds,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Configuration retrieved successfully",
		Data:    configResponse,
	})
}

func toNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}
func toNullInt32(i *int32) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{Int32: *i, Valid: true}
	}
	return sql.NullInt32{}
}
func toNullFloatString(f *float64) sql.NullString {
	if f != nil {
		return sql.NullString{String: strconv.FormatFloat(*f, 'f', -1, 64), Valid: true}
	}
	return sql.NullString{}
}

// POST /api/data/upload
// Description: Client uploads telemetry and test results
// Body: JSON with timestamp, location, signal strength, network type, test results
// Auth: Optional token
func (h *Handler) UploadClientDataHandler(w http.ResponseWriter, r *http.Request) {
	var req MeasurementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}
	if req.Timestamp.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Timestamp is required",
		})
		return
	}
	if req.Latitude == 0 || req.Longitude == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Latitude and longitude are required",
		})
		return
	}
	if req.NetworkType == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Network type is required",
		})
		return
	}
	params := db.CreateMeasurementParams{
		Timestamp:     req.Timestamp,
		Latitude:      strconv.FormatFloat(req.Latitude, 'f', 8, 64),
		Longitude:     strconv.FormatFloat(req.Longitude, 'f', 8, 64),
		NetworkType:   req.NetworkType,
		PlmnID:        toNullString(req.PLMNID),
		Lac:           toNullString(req.LAC),
		Tac:           toNullString(req.TAC),
		CellID:        toNullString(req.CellID),
		Arfcn:         toNullInt32(req.ARFCN),
		FrequencyBand: toNullString(req.FrequencyBand),
		Rssi:          toNullInt32(req.RSSI),
		EcNo:          toNullFloatString(req.EcNo),
		Rscp:          toNullInt32(req.RSCP),
		Rsrq:          toNullFloatString(req.RSRQ),
		Rsrp:          toNullInt32(req.RSRP),
		DownloadSpeed: toNullFloatString(req.DownloadSpeed),
		UploadSpeed:   toNullFloatString(req.UploadSpeed),
		PingMs:        toNullInt32(req.PingMS),
		DnsMs:         toNullInt32(req.DNSMS),
		HttpLatencyMs: toNullInt32(req.HTTPLatencyMS),
		SmsDeliveryMs: toNullInt32(req.SMSDeliveryMS),
	}
	result, err := h.db.CreateMeasurement(r.Context(), params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to store measurement",
		})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to get inserted ID",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Measurement uploaded successfully",
		Data:    MeasurementResponse{ID: id},
	})
}
