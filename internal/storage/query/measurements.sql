-- name: GetClientConfig :one
SELECT id, sampling_interval_seconds, test_types, created_at, updated_at 
FROM client_config
ORDER BY id DESC LIMIT 1;

-- name: UpdateClientConfig :exec
INSERT INTO client_config (
    sampling_interval_seconds,
    test_types
) VALUES (
    ?, ?
);

-- name: CreateNetworkMeasurement :execresult
INSERT INTO network_measurements (
    device_id,
    timestamp,
    latitude,
    longitude,
    cellular_technology,
    network_operator_name,
    network_operator_mccmnc,
    lac,
    rac,
    tac,
    cell_id,
    arfcn,
    signal_strength,
    signal_quality,
    additional_params
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: CreateHttpTestResult :execresult
INSERT INTO http_test_results (
    device_id,
    timestamp,
    success,
    error_message,
    test_type,
    url,
    bytes_transferred,
    duration_millis,
    throughput_kbps
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: CreatePingTestResult :execresult
INSERT INTO ping_test_results (
    device_id,
    timestamp,
    success,
    error_message,
    host,
    packets_transmitted,
    packets_received,
    min_rtt_millis,
    avg_rtt_millis,
    max_rtt_millis,
    ip_address
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: CreateDnsTestResult :execresult
INSERT INTO dns_test_results (
    device_id,
    timestamp,
    success,
    error_message,
    hostname,
    resolution_time_millis,
    resolved_ip_addresses
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
);

-- name: CreateWebTestResult :execresult
INSERT INTO web_test_results (
    device_id,
    timestamp,
    success,
    error_message,
    url,
    loading_time_millis,
    response_code,
    page_size_bytes
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: CreateSmsTestResult :execresult
INSERT INTO sms_test_results (
    device_id,
    timestamp,
    success,
    error_message,
    destination_address,
    submission_time_millis,
    delivery_time_millis,
    duration_millis
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
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