DROP FUNCTION IF EXISTS create_transfer;
CREATE OR REPLACE FUNCTION create_transfer(
    p_from_account_id bigint,
    p_to_account_id bigint,
    p_amount bigint
)
RETURNS TABLE(id BIGINT, from_account_id bigint, to_account_id bigint, amount bigint, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY
    INSERT INTO transfers (from_account_id, to_account_id, amount)
    VALUES (p_from_account_id, p_to_account_id, p_amount)
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
