#!/bin/sh
set -e

echo "$(date) | Waiting for database at ${DB_HOST}:${DB_PORT}..."

while ! pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; do
  echo "$(date) | Waiting for database..."
  sleep 2
done

echo "$(date) | Database is ready!"

echo "$(date) | Running migrations..."
migrate -path /migrations -database "$DB_URL" up

echo "$(date) | Starting application..."
exec "$@"
