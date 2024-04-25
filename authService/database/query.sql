-- Create a table to store JWT tokens
CREATE TABLE jwt_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    token TEXT NOT NULL,
    expiration TIMESTAMP NOT NULL
);
