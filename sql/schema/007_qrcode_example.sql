-- +goose Up

CREATE TABLE QR_CODES (
                          id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
                          store_id INT REFERENCES stores(id) ON DELETE CASCADE,
                          table_no INT NOT NULL,
                          url TEXT UNIQUE NOT NULL,
                          UNIQUE (store_id, table_no)
);

-- +goose Down
DROP TABLE QR_CODES;