DROP FUNCTION IF EXISTS update_account;
CREATE OR REPLACE FUNCTION update_account(account_id BIGINT, new_balance bigint)
RETURNS TABLE(id BIGINT, owner VARCHAR, balance bigint, unit VARCHAR, created_at TIMESTAMP WITH TIME ZONE) AS $$
BEGIN
    RETURN QUERY UPDATE accounts
    SET balance = new_balance
    WHERE accounts.id = account_id
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
