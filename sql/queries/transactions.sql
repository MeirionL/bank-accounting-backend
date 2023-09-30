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
    account_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetTransactions :many
SELECT * FROM transactions
WHERE account_id = $1
ORDER BY traNsaction_time DESC;

-- name: GetTransactionByID :one
SELECT * FROM transactions 
WHERE account_id = $1 AND id = $2;

-- name: GetTransactionsByAccount :many
SELECT * FROM transactions 
WHERE account_id = $1 AND account_number = $2 AND sort_code = $3;

-- name: GetTransactionsByType :many
SELECT * FROM transactions
WHERE account_id = $1 AND type = $2; 

-- name: GetTransactionsWithLimit :many
SELECT * FROM transactions
WHERE account_id = $1
ORDER BY transaction_time DESC
LIMIT $2;

-- name: GetTransactionsByOthersAccount :many
SELECT * FROM transactions
WHERE account_id = $1 AND account_number = $2 AND sort_code = $3;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE account_id = $1 AND id = $2;