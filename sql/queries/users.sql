-- name: GetUserByEmail :one
SELECT * FROM "user" WHERE email=$1;

-- name: GetAllUsers :many
SELECT * FROM "user";

-- name: GetUserBySession :one
SELECT u.* FROM "user" u JOIN "session" s ON
    u.id = s."userId" WHERE s.token = $1;

-- name: GetUserByID :one
SELECT * FROM "user" WHERE id = $1;