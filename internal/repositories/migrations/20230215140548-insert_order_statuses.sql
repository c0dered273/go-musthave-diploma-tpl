
-- +migrate Up
INSERT INTO order_status(name)
VALUES ('NEW'),
       ('PROCESSING'),
       ('INVALID'),
       ('PROCESSED');


-- +migrate Down
TRUNCATE TABLE order_status
    RESTART IDENTITY;