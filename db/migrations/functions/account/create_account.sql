DROP FUNCTION IF EXISTS create_account;
CREATE OR REPLACE FUNCTION create_account(
    p_owner_user_id bigint,
    p_balance bigint,
    p_unit VARCHAR
)
RETURNS TABLE(id BIGINT, owner_user_id bigint, balance bigint, unit VARCHAR, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY
    INSERT INTO accounts (owner_user_id, balance, unit)
    VALUES (p_owner_user_id, p_balance, p_unit)
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
