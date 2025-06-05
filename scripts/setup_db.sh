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


# Insert admin user if not exists (admin@mail.com, 1234)
mysql --ssl=0 -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
INSERT INTO users (email, password_hash, role)
SELECT 'admin@mail.com', '\$2a\$10\$Zg0JsZpUvHW9lqAz88D5seti2Fi5MT6QiZQOLLfNeoCTD4Q4RXyl.', 'admin'
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'admin@mail.com'
);
"

echo "Database setup completed!" 