DROP FUNCTION IF EXISTS create_entry;
CREATE OR REPLACE FUNCTION create_entry(
    p_account_id bigint,
    p_amount bigint
)
RETURNS TABLE(id BIGINT, account_id bigint, amount bigint, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY
    INSERT INTO entries (account_id, amount)
    VALUES (p_account_id, p_amount)
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
