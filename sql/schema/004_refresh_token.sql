-- +goose Up
CREATE TABLE refresh_tokens (
    token text primary key,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    user_id UUID not null references users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP not null,
    revoked_at TIMESTAMP
);

-- +goose Down
DROP TABLE refresh_tokens;