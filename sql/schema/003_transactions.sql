-- +goose Up
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    transaction_time TIMESTAMP NOT NULL,
    type TEXT NOT NULL,
    amount REAL NOT NULL,
    pre_balance REAL NOT NULL,
    post_balance REAL NOT NULL,
    new_account BOOL NOT NULL,
    name TEXT NOT NULL,
    account_number TEXT NOT NULL,
    sort_code TEXT NOT NULL, 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE transactions;