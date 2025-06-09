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

-- name: CountMeasurements :one
SELECT COUNT(*) FROM network_measurements
WHERE timestamp BETWEEN ? AND ?; 

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

-- name: GetAverageDownloadSpeed :one
SELECT CAST(COALESCE(AVG(throughput_kbps), 0) AS DECIMAL(10,2)) as average_speed
FROM http_test_results
WHERE timestamp BETWEEN ? AND ?
AND test_type = 'DOWNLOAD'
AND success = true;

-- name: GetAverageUploadSpeed :one
SELECT CAST(COALESCE(AVG(throughput_kbps), 0) AS DECIMAL(10,2)) as average_speed
FROM http_test_results
WHERE timestamp BETWEEN ? AND ?
AND test_type = 'UPLOAD'
AND success = true;

-- name: GetAverageLatency :one
SELECT CAST(COALESCE(AVG(avg_rtt_millis), 0) AS DECIMAL(10,2)) as average_latency
FROM ping_test_results
WHERE timestamp BETWEEN ? AND ?
AND success = true;


-- name: GetHttpSuccessAndTotal :one
SELECT 
    COUNT(CASE WHEN success = true THEN 1 END) as http_success,
    COUNT(*) as http_total
FROM http_test_results
WHERE timestamp BETWEEN ? AND ?;

-- name: GetPingSuccessAndTotal :one
SELECT 
    COUNT(CASE WHEN success = true THEN 1 END) as ping_success,
    COUNT(*) as ping_total
FROM ping_test_results
WHERE timestamp BETWEEN ? AND ?;

-- name: GetDnsSuccessAndTotal :one
SELECT 
    COUNT(CASE WHEN success = true THEN 1 END) as dns_success,
    COUNT(*) as dns_total
FROM dns_test_results
WHERE timestamp BETWEEN ? AND ?;

-- name: GetWebSuccessAndTotal :one
SELECT 
    COUNT(CASE WHEN success = true THEN 1 END) as web_success,
    COUNT(*) as web_total
FROM web_test_results
WHERE timestamp BETWEEN ? AND ?;

-- name: GetSmsSuccessAndTotal :one
SELECT 
    COUNT(CASE WHEN success = true THEN 1 END) as sms_success,
    COUNT(*) as sms_total
FROM sms_test_results
WHERE timestamp BETWEEN ? AND ?;

-- name: GetNetworkTechnologyDistribution :many
SELECT 
    cellular_technology,
    CAST(COUNT(cellular_technology) * 100.0 / COUNT(*) AS DECIMAL(10,2)) as percentage
FROM network_measurements
WHERE timestamp BETWEEN ? AND ? 
GROUP BY cellular_technology;

-- name: GetMapData :many
SELECT 
    latitude,
    longitude,
    timestamp,
    cellular_technology,
    network_operator_name,
    signal_strength,
    signal_quality,
    additional_params
FROM network_measurements
WHERE 
    latitude BETWEEN ? AND ?
    AND longitude BETWEEN ? AND ?
ORDER BY timestamp DESC;

-- name: ListMeasurements :many
SELECT * FROM network_measurements
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: ListPingMeasurements :many
SELECT * FROM ping_test_results
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: ListHttpMeasurements :many
SELECT * FROM http_test_results
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: ListDnsMeasurements :many
SELECT * FROM dns_test_results
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: ListWebMeasurements :many
SELECT * FROM web_test_results
ORDER BY id DESC
LIMIT ? OFFSET ?;

-- name: ListSmsMeasurements :many
SELECT * FROM sms_test_results
ORDER BY id DESC
LIMIT ? OFFSET ?;
