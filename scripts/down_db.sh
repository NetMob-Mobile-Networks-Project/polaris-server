#!/bin/bash

set -e

# Required environment variables: DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME
: "${DB_HOST:?Need to set DB_HOST}"
: "${DB_PORT:?Need to set DB_PORT}"
: "${DB_USER:?Need to set DB_USER}"
: "${DB_PASS:?Need to set DB_PASS}"
: "${DB_NAME:?Need to set DB_NAME}"

TABLES=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" -Nse "SHOW TABLES FROM \`$DB_NAME\`;")

if [ -z "$TABLES" ]; then
    echo "No tables found in database '$DB_NAME'."
    exit 0
fi

echo "Dropping all tables in database '$DB_NAME'..."

for t in $TABLES; do
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" -e "DROP TABLE \`$DB_NAME\`.\`$t\`;"
done

echo "All tables dropped."
