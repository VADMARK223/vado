#!/bin/sh
set -e

DB_HOST=${DB_HOST:-db}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-vadmark}

echo "⏳ Waiting for postgres at ${DB_HOST}:${DB_PORT}..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; do
echo "Postgres is not ready yet..."
sleep 1
done

echo "✅ Postgres is ready — starting app"
exec "$@"
