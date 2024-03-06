DROP FUNCTION IF EXISTS add_account_balance;
CREATE OR REPLACE FUNCTION add_account_balance(
    p_amount bigint,
    p_id INT
)
RETURNS TABLE(id INT, owner TEXT, balance bigint, unit TEXT, created_at TIMESTAMP)
AS $$
BEGIN
    RETURN QUERY
    UPDATE accounts
    SET balance = balance + p_amount
    WHERE id = p_id
    RETURNING id, owner, balance, unit, created_at;
END;
$$ LANGUAGE plpgsql;
