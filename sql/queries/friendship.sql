-- name: GetFriends :many
SELECT u.* FROM users u JOIN friendship f ON u.id = f.user_id
WHERE u.id = $1;