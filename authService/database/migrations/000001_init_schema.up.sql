CREATE TABLE tokens (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  token TEXT NOT NULL UNIQUE, -- Unique constraint on token
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expiration TIMESTAMP NOT NULL,
  CONSTRAINT idx_jwt_tokens_expiration_user_id UNIQUE (expiration, user_id) -- Composite unique constraint
);