#!/bin/bash

# Load environment variables
source .env

# Create database if it doesn't exist
mysql -u$DB_USER -p$DB_PASS -h$DB_HOST -P$DB_PORT -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;"

# Apply schema
mysql -u$DB_USER -p$DB_PASS -h$DB_HOST -P$DB_PORT $DB_NAME < internal/storage/schema/schema.sql

echo "Database setup completed!" 