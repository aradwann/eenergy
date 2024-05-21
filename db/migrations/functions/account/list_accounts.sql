DROP FUNCTION IF EXISTS list_accounts;
CREATE OR REPLACE FUNCTION list_accounts(p_owner_user_id bigint, limit_value INTEGER, offset_value INTEGER)
RETURNS TABLE(id BIGINT, owner_user_id bigint, balance bigint, unit VARCHAR, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY SELECT a.id, a.owner_user_id, a.balance, a.unit, a.created_at
                 FROM accounts a
                 WHERE a.owner_user_id = owner_user_id
                 ORDER BY a.id
                 LIMIT limit_value
                 OFFSET offset_value;
END;
$$ LANGUAGE plpgsql;
