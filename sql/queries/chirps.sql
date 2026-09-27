-- name: CreateChirps :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    sqlc.arg(body),
    sqlc.arg(user_id)
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;
