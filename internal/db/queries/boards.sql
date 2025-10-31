-- name: GetAllBoards :many
SELECT id, name, description, created_at, updated_at
FROM boards;

-- name: GetAllBoardsByUserId :many
SELECT id, name, description, user_id, created_at, updated_at
FROM boards
WHERE user_id = $1;

-- name: CreateBoard :one
INSERT INTO boards (name, description, user_id)
VALUES ($1, $2, $3)
RETURNING id, name, description, user_id, created_at, updated_at;

-- name: DeleteBoardById :execrows
DELETE FROM boards
WHERE id = $1;

-- name: GetBoardById :one
SELECT id, name, description, created_at, updated_at 
FROM boards
WHERE id = $1;

-- name: UpdateBoardById :one
UPDATE boards
SET name = $2, description = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, description, created_at, updated_at;
