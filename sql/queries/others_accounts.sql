-- name: CreateOthersAccount :exec
INSERT INTO others_accounts (
created_at,
updated_at,
account_name,
account_number,
sort_code,
account_id
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOthersAccountByDetails :one
SELECT * FROM others_accounts WHERE account_id = $1 AND account_number = $2 AND sort_code = $3;

-- name: GetOthersAccountByID :one
SELECT * FROM others_accounts WHERE account_id = $1 AND id = $2;

-- name: GetOthersAccounts :many
SELECT * FROM others_accounts WHERE account_id = $1;

-- name: UpdateOthersAccountName :exec
UPDATE others_accounts
SET account_name = $1
WHERE account_number = $2 AND sort_code = $3 AND account_id = $4;