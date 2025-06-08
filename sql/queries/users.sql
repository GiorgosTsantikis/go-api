-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email=$1;

-- name: GetAllUsers :many
SELECT * FROM "user";