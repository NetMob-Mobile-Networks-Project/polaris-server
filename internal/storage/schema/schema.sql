CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
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
    thresholds JSON NOT NULL, -- Map of threshold values for each metric
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