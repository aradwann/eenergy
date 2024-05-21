DROP FUNCTION IF EXISTS create_account;
CREATE OR REPLACE FUNCTION create_account(
    p_owner VARCHAR,
    p_balance bigint,
    p_unit VARCHAR
)
RETURNS TABLE(id BIGINT, owner VARCHAR, balance bigint, unit VARCHAR, created_at timestamptz) AS $$
BEGIN
    RETURN QUERY
    INSERT INTO accounts (owner, balance, unit)
    VALUES (p_owner, p_balance, p_unit)
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
