
-- +migrate Up
CREATE TABLE IF NOT EXISTS order_status
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(32) UNIQUE NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS order_status;
