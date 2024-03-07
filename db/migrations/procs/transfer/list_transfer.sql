DROP FUNCTION IF EXISTS list_transfers;
CREATE OR REPLACE FUNCTION list_transfers(p_from_account_id BIGINT, p_to_account_id BIGINT, limit_value INTEGER, offset_value INTEGER)
RETURNS TABLE(id BIGINT, from_account_id BIGINT, to_account_id BIGINT, amount bigint, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY SELECT t.id, t.from_account_id, t.to_account_id, t.amount, t.created_at
                 FROM transfers t
                 WHERE t.from_account_id = p_from_account_id OR t.to_account_id = p_to_account_id
                 ORDER BY t.id
                 LIMIT limit_value
                 OFFSET offset_value;
END;
$$ LANGUAGE plpgsql;
