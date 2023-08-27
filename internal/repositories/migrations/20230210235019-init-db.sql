
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

INSERT INTO order_status(name)
VALUES ('NEW'),
       ('PROCESSING'),
       ('INVALID'),
       ('PROCESSED');

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

create or replace procedure withdraw_from_user_balance(
    user_name varchar(50),
    order_id varchar(24),
    amount DECIMAL(16,2)
)
    language plpgsql
as $$
declare
    user_balance users.balance%type := 0;
begin
    select u.balance from users u
    into user_balance
        where u.username = user_name;

    if user_balance < amount then
        raise exception 'balance not enough';
    end if;

    update users
    set balance = balance - amount
    where username = user_name;

    insert into withdrawals(order_id, user_id, amount, processed_at)
    values(order_id,
           (select u.id from users u where u.username = user_name),
           amount,
           now());

    commit;
end; $$
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP EXTENSION IF EXISTS pgcrypto;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS order_status;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS withdrawals;
drop procedure withdraw_from_user_balance(character varying,character varying,numeric);
-- +migrate StatementEnd