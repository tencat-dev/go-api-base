-- name: ExistsPermission :one
SELECT EXISTS(
    SELECT 1 FROM permissions
    WHERE id = $1
);

-- name: InsertPermission :exec
INSERT INTO permissions (object, action)
VALUES ($1, $2)
    ON CONFLICT (object, action) DO NOTHING;