-- +goose Up
CREATE TABLE others_accounts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    account_name TEXT NOT NULL,
    account_number TEXT NOT NULL,
    sort_code TEXT NOT NULL,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE others_accounts;