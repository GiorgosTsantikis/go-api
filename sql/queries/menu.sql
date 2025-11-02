-- name: CreateMenuItem :one
INSERT INTO menu_item ("name", "price") VALUES ($1, $2) RETURNING *;

-- name: CreateConfigOption :one
INSERT INTO config_option ("menu_item_id", "name", "max_selectable") VALUES($1, $2, $3) RETURNING *;

-- name: CreateOption :one
INSERT INTO option ("config_option_id", "name", "price_delta") VALUES($1, $2, $3) RETURNING *;

-- name: GetAllPreviouslyUsedOptions :many
SELECT o.id, o.name, o.price_delta FROM option o INNER JOIN config_option co ON o.config_option_id = co.id
INNER JOIN menu_item mi ON mi.id = $1;

-- name: DeleteMenuItem :exec
DELETE FROM menu_item WHERE ID = $1;
