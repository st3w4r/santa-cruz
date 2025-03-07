-- name: GetDatabase :one
SELECT
    id,
    name,
    path,
    description,
    created_at,
    updated_at
FROM databases
WHERE id = ?;

-- name: ListDatabases :many
SELECT
    id,
    name,
    path,
    description,
    created_at,
    updated_at
FROM databases;

-- name: CreateDatabase :one
INSERT INTO databases (
    name,
    path,
    description,
    created_at,
    updated_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?
) RETURNING *;

-- name: UpdateDatabase :exec
UPDATE databases
SET
    name = ?,
    path = ?,
    description = ?,
    updated_at = ?
WHERE id = ?;

-- name: DeleteDatabase :exec
DELETE FROM databases
WHERE id = ?;
