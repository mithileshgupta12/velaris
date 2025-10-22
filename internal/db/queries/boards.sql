-- name: GetAllBoards :many
SELECT id, name, description, created_at, updated_at
FROM boards;

-- name: CreateBoard :one
INSERT INTO boards (name, description)
VALUES ($1, $2)
RETURNING id, name, description, created_at, updated_at;

-- name: DeleteBoard :execrows
DELETE FROM boards
WHERE id = $1;

-- name: GetBoardById :one
SELECT id, name, description, created_at, updated_at 
FROM boards
WHERE id = $1;
