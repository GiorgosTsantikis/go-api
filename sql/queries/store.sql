-- name: CreateStore :one
INSERT INTO stores (name, tables) VALUES ($1, $2) RETURNING *;

-- name: GenerateQRCode :one
INSERT INTO QR_CODES (store_id, table_no, url) VALUES ($1, $2, $3) RETURNING *;

-- name: GetStoreByID :one
SELECT * FROM stores WHERE id = $1;