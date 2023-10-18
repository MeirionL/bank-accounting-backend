-- +goose Up
CREATE TABLE revoked_tokens (
    token TEXT NOT NULL,
    revoked_at TIMESTAMP NOT NULL,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE revoked_tokens;