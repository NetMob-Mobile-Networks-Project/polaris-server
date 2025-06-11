package api

import (
	"database/sql"
	"encoding/json"
)

// NullStringToString converts sql.NullString to string, returning empty string if null
func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// NullInt32ToInt converts sql.NullInt32 to int, returning 0 if null
func NullInt32ToInt(ni sql.NullInt32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return 0
}

// NullInt64ToInt converts sql.NullInt64 to int64, returning 0 if null
func NullInt64ToInt(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

// NullFloat64ToFloat converts sql.NullFloat64 to float64, returning 0.0 if null
func NullFloat64ToFloat(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0.0
}

// RawMessageToString converts json.RawMessage to string, returning empty string if null
func RawMessageToString(rm json.RawMessage) string {
	if rm == nil {
		return ""
	}
	return string(rm)
}
