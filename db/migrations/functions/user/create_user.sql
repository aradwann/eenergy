DROP FUNCTION IF EXISTS create_user;

CREATE OR REPLACE FUNCTION create_user(
    p_username VARCHAR,
    p_hashed_password VARCHAR,
    p_fullname VARCHAR,
    p_email VARCHAR,
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
DECLARE
    user_role_id BIGINT;
BEGIN
    -- Check if p_role_id is NULL and set it to the 'user' role_id if so
    IF p_role_id IS NULL THEN
        SELECT id INTO user_role_id FROM roles WHERE role_name = 'user';
        
        -- Raise an exception if 'user' role does not exist
        IF user_role_id IS NULL THEN
            RAISE EXCEPTION 'Role "user" does not exist in the roles table';
        END IF;
    ELSE
        user_role_id := p_role_id;
    END IF;

    RETURN QUERY
    INSERT INTO users (username, hashed_password, fullname, email, role_id)
    VALUES (p_username, p_hashed_password, p_fullname, p_email, user_role_id)
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
