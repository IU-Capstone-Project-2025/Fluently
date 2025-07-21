#!/bin/bash

# Database migration script for Fluently backend
# This script runs the migration to fix user and preferences relationship

set -e  # Exit on any error

echo "Starting database migration..."

# Check if we're in the right directory
if [ ! -f "cmd/main.go" ]; then
    echo "Error: Please run this script from the backend directory"
    exit 1
fi

# Check if migration file exists
if [ ! -f "migration_fix.sql" ]; then
    echo "Error: migration_fix.sql not found"
    exit 1
fi

# Get database connection details from environment or config
# You may need to adjust these based on your configuration
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${DB_NAME:-"fluently"}
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-""}

echo "Connecting to database: $DB_HOST:$DB_PORT/$DB_NAME"

# Run the migration
if [ -n "$DB_PASSWORD" ]; then
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f migration_fix.sql
else
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f migration_fix.sql
fi

echo "Migration completed successfully!"
echo "You can now restart your application." 