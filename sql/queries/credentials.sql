-- name: CreateCredentials :one
INSERT INTO credentials (id, username, email, password)
VALUES ($1, $2, $3, $4)
RETURNING *;