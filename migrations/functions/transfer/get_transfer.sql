DROP FUNCTION IF EXISTS get_transfer;
CREATE OR REPLACE FUNCTION get_transfer(p_id BIGINT)
RETURNS TABLE(id BIGINT, from_account_id bigint, to_account_id bigint,amount bigint,  created_at TIMESTAMP WITH TIME ZONE)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.from_account_id, t.to_account_id, t.amount, t.created_at
    FROM transfers t
    WHERE t.id = p_id
    LIMIT 1;
END;
$$;
