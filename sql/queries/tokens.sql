-- name: CreateRevokedToken :exec
INSERT INTO revoked_tokens (token, revoked_at, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRevokedTokens :many
SELECT token FROM revoked_tokens WHERE user_id = $1;