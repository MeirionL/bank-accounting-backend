-- name: CreateAccount :one
INSERT INTO accounts (
    id, 
    created_at, 
    updated_at,  
    account_name, 
    balance,
    account_number, 
    sort_code, 
    user_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAccounts :many
SELECT * FROM accounts WHERE user_id = $1;

-- name: GetAccountByDetails :one
SELECT * FROM accounts WHERE account_number = $1 AND sort_code = $2 AND user_id = $3;

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id = $1 AND user_id = $2;

-- name: UpdateAccount :one
UPDATE accounts
SET updated_at = $3, account_name = $4, balance = $5
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1 AND user_id = $2;

-- name: GetUserIDByAccountID :one
SELECT user_id FROM accounts WHERE id = $1;

-- name: GetAccountsBalances :many
SELECT account_name, balance FROM accounts WHERE user_id = $1;