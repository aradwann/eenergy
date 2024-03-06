DROP FUNCTION IF EXISTS get_account_for_update;
CREATE OR REPLACE FUNCTION get_account_for_update(account_id BIGINT)
RETURNS TABLE(id BIGINT, owner VARCHAR, balance bigint, unit VARCHAR, created_at TIMESTAMP) AS $$
BEGIN
    RETURN QUERY SELECT id, owner, balance, unit, created_at 
                 FROM accounts
                 WHERE id = account_id
                 LIMIT 1
                 FOR NO KEY UPDATE;
END;
$$ LANGUAGE plpgsql;
