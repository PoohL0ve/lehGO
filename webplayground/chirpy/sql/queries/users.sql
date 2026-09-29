-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
    VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- Deletes every user. Irreversible—used by the admin reset endpoint.
-- :exec discards any returned rows; use :execrows if you need the count.
-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;