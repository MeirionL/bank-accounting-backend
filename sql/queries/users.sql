-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, name, hashed_password)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;