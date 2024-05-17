
-- Create a new database named 'authdb'
CREATE DATABASE authdb
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;

-- Connect to the newly created database
\c authdb

-- Create a table within the 'authdb' database
CREATE TABLE tokens (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  token TEXT NOT NULL UNIQUE, -- Unique constraint on token
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expiration TIMESTAMP NOT NULL,
  CONSTRAINT idx_jwt_tokens_expiration_user_id UNIQUE (expiration, user_id) -- Composite unique constraint
);

-- Insert some sample data into the 'users' table
-- INSERT INTO users (username, password) VALUES ('john_doe', 'password123');
-- INSERT INTO users (username, password) VALUES ('jane_smith', 'securepass456');
