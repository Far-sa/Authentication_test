
-- Create a new database named 'authdb'
CREATE DATABASE authdb
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE = template0;

-- Connect to the newly created database
-- \c authdb

-- Create a table within the 'authdb' database
-- CREATE TABLE users (
--     id SERIAL PRIMARY KEY,
--     username VARCHAR(50) NOT NULL,
--     password VARCHAR(50) NOT NULL
-- );

-- Insert some sample data into the 'users' table
-- INSERT INTO users (username, password) VALUES ('john_doe', 'password123');
-- INSERT INTO users (username, password) VALUES ('jane_smith', 'securepass456');
