-- name: CreateTransaction :one
INSERT INTO transactions (
    id, 
    created_at, 
    updated_at, 
    transaction_time, 
    type, 
    amount, 
    pre_balance, 
    post_balance, 
    new_account, 
    name, 
    account_number, 
    sort_code, 
    user_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetTransactions :many
SELECT * FROM transactions;