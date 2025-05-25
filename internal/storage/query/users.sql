-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: CreateUser :execresult
INSERT INTO users (
    username,
    email
) VALUES (
    ?, ?
);

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT ?
OFFSET ?; 