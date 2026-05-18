-- name: InsertRole :exec
INSERT INTO roles (name, description, is_system)
VALUES ($1, $2, $3)
    ON CONFLICT (name) DO NOTHING;

-- name: ExistsRole :one
SELECT EXISTS(
    SELECT 1 FROM roles
    WHERE id = $1
);