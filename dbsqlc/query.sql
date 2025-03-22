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

-- name: GetDatabaseByName :one
SELECT
    id,
    name,
    path,
    description,
    created_at,
    updated_at
FROM databases
WHERE name = ?;

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

-- name: UpdateDatabase :one
UPDATE databases
SET
    name = coalesce(sqlc.narg('name'), name),
    path = coalesce(sqlc.narg('path'), path),
    description = coalesce(sqlc.narg('description'), description),
    updated_at = sqlc.arg('updated_at')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteDatabase :exec
DELETE FROM databases
WHERE id = ?;
