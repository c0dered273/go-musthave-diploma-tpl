
-- +migrate Up
INSERT INTO orders(id, status_id, user_id, amount)
VALUES (00000018,
        (SELECT os.id FROM order_status os WHERE os.name = 'NEW'),
        (SELECT u.id FROM users u WHERE u.username = 'User'),
        999.91),
       (00000026,
        (SELECT os.id FROM order_status os WHERE os.name = 'PROCESSING'),
        (SELECT u.id FROM users u WHERE u.username = 'User'),
        9.00),
       (00000034,
        (SELECT os.id FROM order_status os WHERE os.name = 'PROCESSED'),
        (SELECT u.id FROM users u WHERE u.username = 'User'),
        100.10),
       (00000042,
        (SELECT os.id FROM order_status os WHERE os.name = 'INVALID'),
        (SELECT u.id FROM users u WHERE u.username = 'User'),
        0);

-- +migrate Down
TRUNCATE TABLE orders
    RESTART IDENTITY;