DROP FUNCTION IF EXISTS update_user;

CREATE OR REPLACE FUNCTION update_user(
    p_user_id BIGINT,
    p_username VARCHAR DEFAULT NULL,
    p_hashed_password VARCHAR DEFAULT NULL,
    p_password_changed_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    p_fullname VARCHAR DEFAULT NULL,
    p_email VARCHAR DEFAULT NULL,
    p_is_email_verified BOOLEAN DEFAULT NULL,
    p_role_id BIGINT DEFAULT NULL
)
RETURNS TABLE (
    id BIGINT,
    username VARCHAR,
    hashed_password VARCHAR,
    fullname VARCHAR,
    email VARCHAR,
    password_changed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    is_email_verified BOOLEAN,
    role_id BIGINT
) AS $$
BEGIN
    RETURN QUERY
    UPDATE users
    SET
        username = COALESCE(p_username, users.username),
        hashed_password = COALESCE(p_hashed_password, users.hashed_password),
        password_changed_at = COALESCE(p_password_changed_at, users.password_changed_at),
        fullname = COALESCE(p_fullname, users.fullname),
        email = COALESCE(p_email, users.email),
        is_email_verified = COALESCE(p_is_email_verified, users.is_email_verified),
        role_id = COALESCE(p_role_id, users.role_id)
    WHERE users.id = p_user_id
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
