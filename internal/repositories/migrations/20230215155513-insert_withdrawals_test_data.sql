
-- +migrate Up
INSERT INTO withdrawals(order_id, amount, processed_at)
VALUES (34,
        23.01,
        timestamp '2023-02-15 01:00'),
       (34,
        27.99,
        timestamp '2023-02-15 09:00');

-- +migrate Down
TRUNCATE TABLE withdrawals
    RESTART IDENTITY;