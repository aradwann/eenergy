CREATE OR REPLACE PROCEDURE create_user(
    INOUT p_username VARCHAR,
    INOUT p_hashed_password VARCHAR,
    INOUT p_fullname VARCHAR,
    INOUT p_email VARCHAR,
    OUT password_changed_at_out TIMESTAMP WITH TIME ZONE,
    OUT created_at_out TIMESTAMP WITH TIME ZONE
)
LANGUAGE plpgsql
AS $$
BEGIN
    -- Insert new user and get the password_changed_at and created_at values
    INSERT INTO users (username, hashed_password, fullname, email)
    VALUES (p_username, p_hashed_password, p_fullname, p_email)
    RETURNING password_changed_at, created_at
    INTO password_changed_at_out, created_at_out;

    -- Optionally update the input parameters with the inserted values
    SELECT username, hashed_password, fullname, email
    INTO p_username, p_hashed_password, p_fullname, p_email
    FROM users
    WHERE username = p_username;

END;
$$;
