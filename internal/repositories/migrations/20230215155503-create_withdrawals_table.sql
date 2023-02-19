
-- +migrate Up
CREATE TABLE IF NOT EXISTS withdrawals
(
    id           SERIAL PRIMARY KEY,
    order_id     VARCHAR(24)    NOT NULL,
    user_id      INT            NOT NULL,
    amount       DECIMAL(16, 2) NOT NULL,
    processed_at TIMESTAMP      NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE IF EXISTS withdrawals;