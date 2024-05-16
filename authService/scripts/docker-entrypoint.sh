#!/bin/bash

# Check if database exists
if ! psql -U postgres -h postgres -lqt | grep -q "authDB>"; then
  # Create the database
  psql -U postgres -h postgres -c "CREATE DATABASE authDB"
fi

# Start the Postgres server
exec "$@"
