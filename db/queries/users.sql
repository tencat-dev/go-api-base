-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: InsertUser :one
INSERT INTO users (name, email, password_hash, token_version)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    name = $2,
    email = $3,
    updated_at = now()
WHERE id = $1
    RETURNING *;

-- name: ListUsers :many
SELECT * FROM users LIMIT $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: ExistsUser :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE id = $1
);