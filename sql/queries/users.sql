-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
where email = $1;

-- name: GetUserByID :one
SELECT * FROM users
where id = $1;

-- name: UpdateUser :exec
Update users
set hashed_password = $1, email = $2, updated_at = NOW()
where id = $3;

-- name: DeleteAllUsers :exec
DELETE FROM users;