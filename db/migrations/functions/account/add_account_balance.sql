DROP FUNCTION IF EXISTS add_account_balance;
CREATE OR REPLACE FUNCTION add_account_balance(
    p_amount bigint,
    p_id bigint
)
RETURNS TABLE(id bigint, owner varchar, balance bigint, unit varchar, created_at timestamp with time zone)
AS $$
BEGIN
    RETURN QUERY
    UPDATE accounts 
    SET balance = accounts.balance + p_amount
    WHERE accounts.id = p_id
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
