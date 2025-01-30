-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
Select user_id from refresh_tokens
where (token = $1) and (expires_at > NOW()) and (revoked_at is NULL);

-- name: UpdateRefreshToken :exec
Update refresh_tokens
set updated_at = NOW(), revoked_at = NOW()
where token = $1;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM refresh_tokens;