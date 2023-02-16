
-- +migrate Up
CREATE TABLE IF NOT EXISTS orders
(
    id        DECIMAL PRIMARY KEY,
    status_id INT NOT NULL,
    user_id   INT         NOT NULL,
    amount    DECIMAL(16, 2),
    uploaded_at TIMESTAMP,
    CONSTRAINT fk_status FOREIGN KEY (status_id) REFERENCES order_status (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);


-- +migrate Down
DROP TABLE IF EXISTS orders;