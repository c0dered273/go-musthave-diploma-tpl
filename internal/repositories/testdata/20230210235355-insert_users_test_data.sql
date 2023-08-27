
-- +migrate Up
INSERT INTO users(username, password, balance)
VALUES ('Admin', crypt('admin', gen_salt('bf')), 0),
       ('User', crypt('user', gen_salt('bf')), 30) ;

-- +migrate Down
TRUNCATE TABLE users
    RESTART IDENTITY;