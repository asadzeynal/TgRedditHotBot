#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$TGRHB_DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
