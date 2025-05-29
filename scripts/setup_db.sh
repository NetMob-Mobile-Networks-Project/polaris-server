#!/bin/bash

echo "Waiting for MySQL on $DB_HOST:$DB_PORT ..."
until mysqladmin ping --ssl=0 -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" --silent 2>/dev/null; do
    sleep 2
done
echo "MySQL is up. Initializing schema"

# Create database if it doesn't exist
mysql --ssl=0 -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;"

# Apply schema
mysql --ssl=0 -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < internal/storage/schema/schema.sql

echo "Database setup completed!" 