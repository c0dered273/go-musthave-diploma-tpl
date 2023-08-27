
-- +migrate Up
INSERT INTO withdrawals(order_id, user_id, amount, processed_at)
VALUES ('3434',
        (SELECT u.id FROM users u WHERE username = 'User'),
        23.00,
        timestamp '2023-02-15 01:00'),
       ('4242',
        (SELECT u.id FROM users u WHERE username = 'User'),
        27.00,
        timestamp '2023-02-15 09:00');

-- +migrate Down
TRUNCATE TABLE withdrawals
    RESTART IDENTITY;