#!/bin/sh

echo "Running migrations"

DBSTRING="postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE"

migrate -database "${DBSTRING}" -path /migrations $1

echo "Successfully run migrations"