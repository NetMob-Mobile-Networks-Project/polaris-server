-- name: GetClientConfig :one
SELECT * FROM client_config
ORDER BY id DESC LIMIT 1;

-- name: UpdateClientConfig :exec
INSERT INTO client_config (
    sampling_interval_seconds,
    test_types,
    thresholds
) VALUES (
    ?, ?, ?
);

-- name: CreateMeasurement :execresult
INSERT INTO measurements (
    timestamp,
    latitude,
    longitude,
    network_type,
    plmn_id,
    lac,
    tac,
    cell_id,
    arfcn,
    frequency_band,
    rssi,
    ec_no,
    rscp,
    rsrq,
    rsrp,
    download_speed,
    upload_speed,
    ping_ms,
    dns_ms,
    http_latency_ms,
    sms_delivery_ms
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetMeasurement :one
SELECT * FROM measurements
WHERE id = ? LIMIT 1;

-- name: ListMeasurements :many
SELECT * FROM measurements
WHERE timestamp BETWEEN ? AND ?
ORDER BY timestamp DESC
LIMIT ? OFFSET ?;

-- name: CountMeasurements :one
SELECT COUNT(*) FROM measurements
WHERE timestamp BETWEEN ? AND ?; 