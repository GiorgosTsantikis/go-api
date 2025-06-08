-- name: GetFriends :many
SELECT u.* FROM "user" u JOIN friendship f ON u.id = f.user_id
WHERE u.id = $1;