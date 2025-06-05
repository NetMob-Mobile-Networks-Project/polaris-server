CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('admin', 'user') NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS client_config (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sampling_interval_seconds INT NOT NULL DEFAULT 300, -- 5 minutes default
    test_types JSON NOT NULL, -- Array of test types to run
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS measurements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    network_type VARCHAR(50) NOT NULL, -- 2G, 3G, 4G, 5G, etc.
    plmn_id VARCHAR(10), -- Public Land Mobile Network ID
    lac VARCHAR(10), -- Location Area Code
    tac VARCHAR(10), -- Tracking Area Code
    cell_id VARCHAR(20), -- Cell ID
    arfcn INT, -- Absolute Radio Frequency Channel Number
    frequency_band VARCHAR(20), -- Frequency band (e.g., "B1", "B3")
    rssi INT, -- Received Signal Strength Indicator (dBm)
    ec_no DECIMAL(5,2), -- Energy per chip to noise power spectral density ratio (dB)
    rscp INT, -- Received Signal Code Power (dBm)
    rsrq DECIMAL(5,2), -- Reference Signal Received Quality (dB)
    rsrp INT, -- Reference Signal Received Power (dBm)
    download_speed DECIMAL(10,2), -- Mbps
    upload_speed DECIMAL(10,2), -- Mbps
    ping_ms INT, -- Milliseconds
    dns_ms INT, -- DNS resolution time in milliseconds
    http_latency_ms INT, -- HTTP GET latency in milliseconds
    sms_delivery_ms INT, -- SMS delivery time in milliseconds
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_timestamp (timestamp),
    INDEX idx_location (latitude, longitude),
    INDEX idx_network (network_type)
);

CREATE TABLE IF NOT EXISTS network_measurements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    cellular_technology VARCHAR(50),
    network_operator_name VARCHAR(255),
    network_operator_mccmnc VARCHAR(10),
    lac INT,
    rac INT,
    tac INT,
    cell_id BIGINT,
    arfcn INT,
    signal_strength INT,
    signal_quality INT,
    additional_params JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_device_timestamp (device_id, timestamp),
    INDEX idx_location (latitude, longitude)
);

CREATE TABLE IF NOT EXISTS http_test_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    test_type ENUM('DOWNLOAD', 'UPLOAD') NOT NULL,
    url VARCHAR(2048) NOT NULL,
    bytes_transferred BIGINT,
    duration_millis BIGINT,
    throughput_kbps DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_device_timestamp (device_id, timestamp)
);

CREATE TABLE IF NOT EXISTS ping_test_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    host VARCHAR(255) NOT NULL,
    packets_transmitted INT,
    packets_received INT,
    min_rtt_millis FLOAT,
    avg_rtt_millis FLOAT,
    max_rtt_millis FLOAT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_device_timestamp (device_id, timestamp)
);

CREATE TABLE IF NOT EXISTS dns_test_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    hostname VARCHAR(255) NOT NULL,
    resolution_time_millis BIGINT,
    resolved_ip_addresses JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_device_timestamp (device_id, timestamp)
);

CREATE TABLE IF NOT EXISTS web_test_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    url VARCHAR(2048) NOT NULL,
    loading_time_millis BIGINT,
    response_code INT,
    page_size_bytes BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_device_timestamp (device_id, timestamp)
);

CREATE TABLE IF NOT EXISTS sms_test_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    destination_address VARCHAR(20) NOT NULL,
    submission_time_millis BIGINT,
    delivery_time_millis BIGINT,
    duration_millis BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_device_timestamp (device_id, timestamp)
); 