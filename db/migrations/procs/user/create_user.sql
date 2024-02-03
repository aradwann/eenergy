DROP FUNCTION IF EXISTS create_user;
CREATE OR REPLACE FUNCTION create_user(
    p_username VARCHAR,
    p_hashed_password VARCHAR,
    p_fullname VARCHAR,
    p_email VARCHAR
)
RETURNS TABLE (
    username VARCHAR,
    hashed_password VARCHAR,
    fullname VARCHAR,
    email VARCHAR,
    password_changed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    is_email_verified BOOLEAN,
    role VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    INSERT INTO users (username, hashed_password, fullname, email)
    VALUES (p_username, p_hashed_password, p_fullname, p_email)
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
