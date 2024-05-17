-- Create the database
CREATE DATABASE mysqlUser
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

-- Grant privileges to the user (assuming 'root' user is creating and granting privileges to 'root' user)
-- Normally, you would create a new user and grant privileges to that user, but this step is optional
GRANT ALL PRIVILEGES ON mysqlUser.* TO 'root'@'%' IDENTIFIED BY 'password';
FLUSH PRIVILEGES;

-- Connect to the newly created database
USE mysqlUser;

-- Create the 'users' table within the 'mysqlUser' database
CREATE TABLE users(
    id INT AUTO_INCREMENT PRIMARY KEY,
    phone_number VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Verify the table creation
SHOW TABLES;

-- Insert some sample data into the 'users' table
-- INSERT INTO users (username, password) VALUES ('john_doe', 'password123');
-- INSERT INTO users (username, password) VALUES ('jane_smith', 'securepass456');
