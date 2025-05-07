#!/bin/sh
set -e

echo "Running database migrations..."
sql-migrate up -config=dbconfig.yml -env=development

echo "Starting movies_service..."
exec "/app/movies_service"
