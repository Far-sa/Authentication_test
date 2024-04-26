-- Create a table to store JWT tokens
CREATE TABLE jwt_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token TEXT NOT NULL,
    expiration TIMESTAMP NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);