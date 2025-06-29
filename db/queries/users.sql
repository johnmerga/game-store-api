-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, first_name, last_name, role, phone
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users 
SET first_name = $2, last_name = $3, phone = $4, avatar_url = $5
WHERE id = $1 
RETURNING *;

-- name: UpdateUserStatus :exec
UPDATE users SET status = $2 WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users 
WHERE ($1::user_role IS NULL OR role = $1)
AND ($2::user_status IS NULL OR status = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
