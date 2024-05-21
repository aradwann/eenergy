DROP FUNCTION IF EXISTS list_entries;
CREATE OR REPLACE FUNCTION list_entries(p_account_id BIGINT, limit_value INTEGER, offset_value INTEGER)
RETURNS TABLE(id BIGINT, account_id BIGINT, amount bigint, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY SELECT e.id, e.account_id, e.amount, e.created_at
                 FROM entries e
                 WHERE e.account_id = p_account_id 
                 ORDER BY e.id
                 LIMIT limit_value
                 OFFSET offset_value;
END;
$$ LANGUAGE plpgsql;
