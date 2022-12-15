#!/bin/sh

# We use the set -e command the make sure that the script will exit immediately, if the command returns a non-zero status
set -e

# The first step is to run db migration
echo "run db migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

# Here, we start the app
echo "start the app"
exec "$@"