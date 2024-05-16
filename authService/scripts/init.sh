#!/bin/bash

# Wait for PostgreSQL to be ready
until pg_isready -U postgres -h localhost; do
    >&2 echo "PostgreSQL is not yet ready..."
    sleep 1
done

# Check if the database already exists
if ! psql -U postgres -h localhost -lqt | cut -d \| -f 1 | grep -qw "authDB"; then
    # Create the database
    psql -U postgres -h localhost -c "CREATE DATABASE authDB"
    echo "Database 'authB' created successfully"
else
    echo "Database 'authDB' already exists"
fi
