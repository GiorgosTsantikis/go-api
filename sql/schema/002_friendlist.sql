-- +goose Up
CREATE TYPE friend_request_status AS ENUM ('ACCEPTED', 'REJECTED', 'PENDING');

CREATE TABLE friendship (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    friend_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (friend_id) REFERENCES "user"(id) ON DELETE CASCADE,
    request_status friend_request_status NOT NULL DEFAULT 'PENDING',
    CONSTRAINT unique_friendship UNIQUE ( user_id, friend_id)
);

-- +goose Down
DROP TABLE IF EXISTS friendship;
DROP TYPE IF EXISTS friend_request_status;