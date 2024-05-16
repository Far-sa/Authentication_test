#!/bin/bash

# Function to wait for Postgres to be ready
wait_for_postgres() {
    until pg_isready -U postgres; do
        >&2 echo "Postgres is not yet ready..."
        sleep 1
    done
}

# Function to create the database if it doesn't exist
create_database_if_not_exists() {
    if ! psql -U postgres -h localhost -lqt | grep -q "authDB>"; then
        # Create the database
        #PGPASSWORD="$POSTGRES_PASSWORD" psql -U "$POSTGRES_USER" -h localhost -c "CREATE DATABASE authDB"
        psql -U postgres -h localhost -c "CREATE DATABASE authDB"
    fi
}

# Main function
main() {
     # Start PostgreSQL service
    service postgresql start
    
    wait_for_postgres
    create_database_if_not_exists
    # Start the Postgres server
    exec "$@"
}

main "$@"
