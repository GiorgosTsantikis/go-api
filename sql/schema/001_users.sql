-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT pg_catalog.gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    profilePic TEXT
);
-- +goose Down
DROP TABLE users;