-- name: CreateUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id, name, email, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE id = $1;
