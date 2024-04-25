-- user.sql

-- name: CreateUser
INSERT INTO users (username, email, created_at, updated_at) VALUES (?, ?, NOW(), NOW());

-- name: GetUserByID
SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?;

-- name: UpdateUserByID
UPDATE users SET username = ?, email = ?, updated_at = NOW() WHERE id = ?;

-- name: DeleteUserByID
DELETE FROM users WHERE id = ?;
