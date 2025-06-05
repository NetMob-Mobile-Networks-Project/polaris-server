-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (
    email,
    password_hash,
    role
) VALUES (
    ?, ?, ?
);
-- Follow-up query (if your code supports running multiple statements)
SELECT * FROM users WHERE id = LAST_INSERT_ID();

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, email, role, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateUserExec :exec
UPDATE users
SET 
    email = COALESCE(?, email),
    password_hash = COALESCE(?, password_hash),
    role = COALESCE(?, role)
WHERE id = ?;

-- -- name: GetUserByID :one
-- SELECT * FROM users WHERE id = ?;


-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users; 