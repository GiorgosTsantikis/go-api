-- +goose Up
CREATE TABLE credentials(
    id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE ,
    username VARCHAR(255) NOT NULL UNIQUE ,
    email VARCHAR(255) NOT NULL UNIQUE ,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE credentials;