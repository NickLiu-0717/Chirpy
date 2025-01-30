-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
Select * from chirps
order by created_at;

-- name: GetAuthorChirps :many
Select * from chirps
where user_id = $1
order by created_at;

-- name: GetOneChirp :one
select * from chirps
where id = $1;

-- name: GetUserFromChirp :one
Select user_id from chirps
where id = $1;

-- name: DeleteOneChirp :one
DELETE FROM chirps
WHERE id = $1
RETURNING *;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;