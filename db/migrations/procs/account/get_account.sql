DROP FUNCTION IF EXISTS get_account;
CREATE OR REPLACE FUNCTION get_account(p_id BIGINT)
RETURNS TABLE(id BIGINT, owner VARCHAR, balance bigint, unit VARCHAR, created_at TIMESTAMP WITH TIME ZONE)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT a.id, a.owner, a.balance, a.unit, a.created_at
    FROM accounts a
    WHERE a.id = p_id
    LIMIT 1;
END;
$$;
