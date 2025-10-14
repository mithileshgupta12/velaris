-- name: GetAllBoards :many
SELECT id, name, description, created_at, updated_at
FROM boards;
