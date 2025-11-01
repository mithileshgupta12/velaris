-- name: GetAllBoardsByUserId :many
SELECT id, name, description, user_id, created_at, updated_at
FROM boards
WHERE user_id = $1;

-- name: CreateBoard :one
INSERT INTO boards (name, description, user_id)
VALUES ($1, $2, $3)
RETURNING id, name, description, user_id, created_at, updated_at;

-- name: DeleteBoardByIdAndUserId :execrows
DELETE FROM boards
WHERE id = $1
AND user_id = $2;

-- name: GetBoardByIdAndUserId :one
SELECT id, name, description, user_id, created_at, updated_at 
FROM boards
WHERE id = $1
AND user_id = $2;

-- name: UpdateBoardByIdAndUserId :one
UPDATE boards
SET name = $3, description = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
AND user_id = $2
RETURNING id, name, description, user_id, created_at, updated_at;
