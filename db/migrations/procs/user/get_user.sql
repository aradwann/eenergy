DROP FUNCTION IF EXISTS get_user;
CREATE OR REPLACE FUNCTION get_user(
    p_username VARCHAR
)
RETURNS TABLE (
    username VARCHAR,
    hashed_password VARCHAR,
    fullname VARCHAR,
    email VARCHAR,
    password_changed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        users.username,
        users.hashed_password,
        users.fullname,
        users.email,
        users.password_changed_at,
        users.created_at
    FROM users
    WHERE users.username = p_username
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;