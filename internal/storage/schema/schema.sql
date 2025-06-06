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