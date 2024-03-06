DROP FUNCTION IF EXISTS list_accounts;
CREATE OR REPLACE FUNCTION list_accounts(owner_value VARCHAR, limit_value INTEGER, offset_value INTEGER)
RETURNS TABLE(id BIGINT, owner VARCHAR, balance bigint, unit VARCHAR, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY SELECT a.id, a.owner, a.balance, a.unit, a.created_at
                 FROM accounts a
                 WHERE a.owner = owner_value
                 ORDER BY a.id
                 LIMIT limit_value
                 OFFSET offset_value;
END;
$$ LANGUAGE plpgsql;
