#!/bin/sh

set -e

echo "run db migration"
/usr/local/bin/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"