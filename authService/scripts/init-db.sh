#!/bin/bash

#!/bin/bash
set -e  # Exit immediately if a command exits with a non-zero status

echo "Running init-db.sh"
touch /docker-entrypoint-initdb.d/test_file
ls -l /docker-entrypoint-initdb.d/

# Function to wait for Postgres to be ready
wait_for_postgres() {
    until pg_isready -U postgres -h postgres; do
        >&2 echo "Postgres is not yet ready..."
        sleep 1
    done
}

# Function to create the database if it doesn't exist
create_database_if_not_exists() {
    if ! psql -U postgres -h postgres -lqt | grep -qw "authdb"; then
        # Create the database
        psql -U postgres -h postgres -c "CREATE DATABASE authdb"
    fi
}

# Main function
main() {
    
    wait_for_postgres
    create_database_if_not_exists
    
    # Start the Postgres server (this line may need to be adjusted based on your setup)
    exec "$@"
}

main "$@"



# set -e

# psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "authDB" <<-EOSQL
#    -- Ensure the dblink extension is installed
#    CREATE EXTENSION IF NOT EXISTS dblink;

#    -- Conditionally create the database
#    DO \$\$
#    BEGIN
#       IF NOT EXISTS (
#          SELECT FROM pg_catalog.pg_database
#          WHERE datname = 'authDB'
#       ) THEN
#          PERFORM dblink_exec('dbname=' || current_database(), 'CREATE DATABASE authDB');
#       END IF;
#    END
#    \$\$;
# EOSQL
