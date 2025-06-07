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
	// Set headers for CSV download
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=telemetry_data.csv")

	// Create CSV writer
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write CSV header
	headers := []string{
		"Timestamp", "Device ID", "Test Type", "Success", "Latitude", "Longitude",
		"Cellular Technology", "Network Operator", "MCCMNC", "LAC", "RAC", "TAC",
		"Cell ID", "ARFCN", "Signal Strength", "Signal Quality",
		"Test Details", "Additional Parameters",
	}
	if err := csvWriter.Write(headers); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to write CSV header")
		return
	}

	// Query all measurements
	measurements, err := h.db.ListMeasurements(r.Context(), db.ListMeasurementsParams{
		FromTimestamp: time.Now().AddDate(0, 0, -30), // Last 30 days by default
		ToTimestamp:   time.Now(),
		Limit:         10000, // Adjust based on your needs
		Offset:        0,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to query measurements")
		return
	}

	// Write each measurement as a CSV row
	for _, m := range measurements {
		// Convert additional parameters to string
		additionalParams := ""
		if m.AdditionalParams != nil {
			additionalParams = string(m.AdditionalParams)
		}

		// Format network metadata
		networkMeta := fmt.Sprintf("Tech: %s, Operator: %s, MCCMNC: %s",
			getStringValue(m.CellularTechnology),
			getStringValue(m.NetworkOperatorName),
			getStringValue(m.NetworkOperatorMccmnc),
		)

		// Create CSV row
		row := []string{
			m.Timestamp.Format(time.RFC3339),
			m.DeviceID,
			"Network Measurement", // Test type
			"true",                // Success
			getStringValue(m.Latitude),
			getStringValue(m.Longitude),
			getStringValue(m.CellularTechnology),
			getStringValue(m.NetworkOperatorName),
			getStringValue(m.NetworkOperatorMccmnc),
			getInt32Value(m.Lac),
			getInt32Value(m.Rac),
			getInt32Value(m.Tac),
			getInt64Value(m.CellID),
			getInt32Value(m.Arfcn),
			getInt32Value(m.SignalStrength),
			getInt32Value(m.SignalQuality),
			networkMeta,
			additionalParams,
		}

		if err := csvWriter.Write(row); err != nil {
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
		FromTimestamp: time.Now().AddDate(0, 0, -30), // Last 30 days by default
		ToTimestamp:   time.Now(),
		Limit:         10000, // Adjust based on your needs
		Offset:        0,
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
