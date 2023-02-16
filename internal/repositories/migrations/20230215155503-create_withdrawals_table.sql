
-- +migrate Up
CREATE TABLE IF NOT EXISTS withdrawals
(
    id SERIAL PRIMARY KEY,
    order_id DECIMAL NOT NULL,
    amount DECIMAL(16,2) NOT NULL ,
    processed_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders (id)
);

-- +migrate Down
DROP TABLE IF EXISTS withdrawals;