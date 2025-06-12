package api

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage/db"
)

// GET /api/export/csv
// Description: Export all telemetry data as CSV
func (h *Handler) ExportCSVHandler(w http.ResponseWriter, r *http.Request) {
	metric := r.URL.Query().Get("metric")

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

	// Get all data without pagination
	limit := int32(1000000) // Large number to get all records
	pageNum := int32(1)

	if metric == "network" {
		list, err := h.db.ListMeasurements(r.Context(), db.ListMeasurementsParams{
			Limit:         limit,
			Offset:        (pageNum - 1) * limit,
			FromTimestamp: time.Time{}, // Get all time
			ToTimestamp:   time.Now(),
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to export network measurements: %v", err))
			return
		}

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

		for _, m := range list {
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
			FromTimestamp: time.Time{}, // Get all time
			ToTimestamp:   time.Now(),
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to export ping measurements: %v", err))
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
			FromTimestamp: time.Time{}, // Get all time
			ToTimestamp:   time.Now(),
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to export http measurements")
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
			FromTimestamp: time.Time{}, // Get all time
			ToTimestamp:   time.Now(),
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to export dns measurements: %v", err))
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
			FromTimestamp: time.Time{}, // Get all time
			ToTimestamp:   time.Now(),
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to export web measurements")
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
			FromTimestamp: time.Time{}, // Get all time
			ToTimestamp:   time.Now(),
		})
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to export sms measurements")
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
			})
		}
	}

	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_measurements.csv", metric))

	// Create CSV writer
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write header
	if err := csvWriter.Write(labels); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to write CSV header")
		return
	}

	// Write data rows
	for _, row := range values {
		var csvRow []string
		for _, label := range labels {
			value := row[label]
			var strValue string
			switch v := value.(type) {
			case time.Time:
				strValue = v.Format(time.RFC3339)
			case nil:
				strValue = ""
			default:
				strValue = fmt.Sprintf("%v", v)
			}
			csvRow = append(csvRow, strValue)
		}
		if err := csvWriter.Write(csvRow); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Failed to write CSV row")
			return
		}
	}
}

// GET /api/export/kml
// Description: Export geo-locations and metadata as KML (for Google Earth)
func (h *Handler) ExportKMLHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for KML download
	w.Header().Set("Content-Type", "application/vnd.google-earth.kml+xml")
	w.Header().Set("Content-Disposition", "attachment; filename=telemetry_data.kml")

	// Write KML header
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
<Document>
<name>Network Measurements</name>
<description>Network measurement data exported from Polaris Server</description>
<Style id="signal_strength">
<IconStyle>
<Icon>
<href>http://maps.google.com/mapfiles/kml/shapes/placemark_circle.png</href>
</Icon>
<color>ff0000ff</color>
<scale>0.8</scale>
</IconStyle>
</Style>
`)

	// Query all measurements
	measurements, err := h.db.ListMeasurements(r.Context(), db.ListMeasurementsParams{
		Limit:         10000,
		Offset:        0,
		FromTimestamp: time.Now().Add(-100 * 24 * time.Hour),
		ToTimestamp:   time.Now(),
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to query measurements")
		return
	}

	// Write each measurement as a KML Placemark
	for _, m := range measurements {
		if !m.Latitude.Valid || !m.Longitude.Valid {
			continue
		}

		// Create description with metadata
		description := fmt.Sprintf(`
			Timestamp: %s
			Device ID: %s
			Cellular Technology: %s
			Network Operator: %s
			Signal Strength: %s dBm
			Signal Quality: %s
			Cell ID: %s
			ARFCN: %s
		`,
			m.Timestamp.Format(time.RFC3339),
			m.DeviceID,
			getStringValue(m.CellularTechnology),
			getStringValue(m.NetworkOperatorName),
			getInt32Value(m.SignalStrength),
			getInt32Value(m.SignalQuality),
			getInt64Value(m.CellID),
			getInt32Value(m.Arfcn),
		)

		// Create name based on signal strength
		name := "Unknown Signal"
		if m.SignalStrength.Valid {
			strength := m.SignalStrength.Int32
			switch {
			case strength >= -70:
				name = "Excellent Signal"
			case strength >= -85:
				name = "Good Signal"
			case strength >= -100:
				name = "Fair Signal"
			default:
				name = "Poor Signal"
			}
		}

		// Write Placemark
		fmt.Fprintf(w, `
<Placemark>
<name>%s</name>
<description><![CDATA[%s]]></description>
<styleUrl>#signal_strength</styleUrl>
<Point>
<coordinates>%s,%s,0</coordinates>
</Point>
</Placemark>
`,
			name,
			description,
			m.Longitude.String,
			m.Latitude.String,
		)
	}

	// Write KML footer
	fmt.Fprintf(w, `
</Document>
</kml>
`)
}

// Helper functions for handling nullable fields
func getStringValue(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func getInt32Value(i sql.NullInt32) string {
	if i.Valid {
		return strconv.Itoa(int(i.Int32))
	}
	return ""
}

func getInt64Value(i sql.NullInt64) string {
	if i.Valid {
		return strconv.FormatInt(i.Int64, 10)
	}
	return ""
}
