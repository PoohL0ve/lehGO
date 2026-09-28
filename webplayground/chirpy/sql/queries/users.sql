-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email)
    VALUES (gen_random_uuid(), NOW(), NOW(), $1)
RETURNING *;

-- Deletes every user. Irreversible—used by the admin reset endpoint.
-- :exec discards any returned rows; use :execrows if you need the count.
-- name: DeleteUsers :exec
DELETE FROM users;