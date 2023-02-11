
-- +migrate Up
INSERT INTO users(username, password)
VALUES ('Admin', crypt('admin', gen_salt('bf'))),
       ('User', crypt('user', gen_salt('bf')));

-- +migrate Down
TRUNCATE TABLE users
    RESTART IDENTITY;