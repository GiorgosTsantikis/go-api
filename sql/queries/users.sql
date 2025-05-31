-- name: CreateUser :one
INSERT INTO users (id, username, profilePic)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByUserName :one
SELECT * FROM users WHERE username=$1;