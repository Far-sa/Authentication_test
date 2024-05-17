-- CREATE DATABASE IF NOT EXISTS authDB;



DO $$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_database
      WHERE datname = 'authDB'
   ) THEN
      PERFORM dblink_exec('dbname=' || current_database(), 'CREATE DATABASE authDB');
   END IF;
END
$$;
