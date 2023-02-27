
-- +migrate Up
-- +migrate StatementBegin
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
drop procedure withdraw_from_user_balance(character varying,character varying,numeric);