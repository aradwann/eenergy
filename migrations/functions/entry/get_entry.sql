DROP FUNCTION IF EXISTS get_entry;
CREATE OR REPLACE FUNCTION get_entry(p_id BIGINT)
RETURNS TABLE(id BIGINT, account_id bigint, amount bigint, created_at TIMESTAMP WITH TIME ZONE)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT e.id,  e.account_id, e.amount, e.created_at
    FROM entries e
    WHERE e.id = p_id
    LIMIT 1;
END;
$$;
