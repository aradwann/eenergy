CREATE OR REPLACE FUNCTION update_user(
    p_username VARCHAR,
    p_hashed_password VARCHAR,
    p_password_changed_at TIMESTAMP WITH TIME ZONE,
    p_fullname VARCHAR,
    p_email VARCHAR
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
    UPDATE users
    SET
        hashed_password = COALESCE(p_hashed_password, users.hashed_password),
        password_changed_at = COALESCE(p_password_changed_at, users.password_changed_at),
        fullname = COALESCE(p_fullname, users.fullname),
        email = COALESCE(p_email, users.email)
    WHERE users.username = p_username
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
