
-- +migrate Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users(
                                    user_id serial PRIMARY KEY,
                                    username VARCHAR(50) UNIQUE NOT NULL,
                                    password VARCHAR(72) NOT NULL
);

-- +migrate Down
DROP EXTENSION IF EXISTS pgcrypto;

DROP TABLE IF EXISTS users;