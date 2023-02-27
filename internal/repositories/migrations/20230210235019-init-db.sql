
-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users
(
    id       serial PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(72)        NOT NULL,
    balance  DECIMAL(16, 2)     NOT NULL
);

CREATE TABLE IF NOT EXISTS order_status
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(32) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS orders
(
    id          DECIMAL PRIMARY KEY,
    status_id   INT NOT NULL,
    user_id     INT NOT NULL,
    amount      DECIMAL(16, 2),
    uploaded_at TIMESTAMP NOT NULL ,
    CONSTRAINT fk_status FOREIGN KEY (status_id) REFERENCES order_status (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS withdrawals
(
    id           SERIAL PRIMARY KEY,
    order_id     VARCHAR(24) UNIQUE NOT NULL,
    user_id      INT                NOT NULL,
    amount       DECIMAL(16, 2)     NOT NULL,
    processed_at TIMESTAMP          NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE OR REPLACE PROCEDURE withdraw_from_user_balance(
    user_name VARCHAR(50),
    orderId VARCHAR(24),
    amount DECIMAL(16, 2)
)
    LANGUAGE plpgsql
AS
$$
DECLARE
    user_balance users.balance%type := 0;
BEGIN
    SELECT u.balance FROM users u INTO user_balance WHERE u.username = user_name;

    IF user_balance < amount THEN
        RAISE EXCEPTION 'balance not enough';
    END IF;

    UPDATE users
    SET balance = balance - amount
    WHERE username = user_name;

    INSERT INTO withdrawals(order_id, user_id, amount, processed_at)
    VALUES (orderId,
            (SELECT u.id FROM users u WHERE u.username = user_name),
            amount,
            now());

    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE update_order_status_and_user_balance(
    user_name VARCHAR(50),
    orderId VARCHAR(24),
    order_status VARCHAR(32),
    incoming_amount DECIMAL(16, 2)
)
    LANGUAGE plpgsql
AS
$$
    DECLARE
        diag_count INT := 0;
BEGIN
    UPDATE orders
    SET status_id = (SELECT os.id FROM order_status os WHERE os.name = order_status),
        amount    = incoming_amount
    WHERE ID = orderId;
    GET DIAGNOSTICS diag_count = ROW_COUNT;

    UPDATE users
    SET balance = balance + incoming_amount
    WHERE username = user_name;
    GET DIAGNOSTICS diag_count = ROW_COUNT;

    COMMIT;
END;
$$;
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP EXTENSION IF EXISTS pgcrypto;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS order_status;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS withdrawals;
DROP PROCEDURE withdraw_from_user_balance(character varying,character varying,numeric);
DROP PROCEDURE update_order_status_and_user_balance(character varying,character varying,character varying,numeric)
-- +migrate StatementEnd