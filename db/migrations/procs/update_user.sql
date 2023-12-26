DROP PROCEDURE IF EXISTS update_user;
CREATE OR REPLACE PROCEDURE update_user(
    p_username VARCHAR,
    p_hashed_password VARCHAR,
    p_password_changed_at TIMESTAMP WITH TIME ZONE,
    p_fullname VARCHAR,
    p_email VARCHAR,
    OUT username_out VARCHAR,
    OUT hashed_password_out VARCHAR,
    OUT fullname_out VARCHAR,
    OUT email_out VARCHAR,
    OUT password_changed_at_out TIMESTAMP WITH TIME ZONE,
    OUT created_at_out TIMESTAMP WITH TIME ZONE
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE users
    SET
        hashed_password = COALESCE(p_hashed_password, users.hashed_password),
        password_changed_at = COALESCE(p_password_changed_at, users.password_changed_at),
        fullname = COALESCE(p_fullname, users.fullname),
        email = COALESCE(p_email, users.email)
    WHERE users.username = p_username
    RETURNING
        username,
        hashed_password,
        fullname,
        email,
        password_changed_at,
        created_at
    INTO
        username_out,
        hashed_password_out,
        fullname_out,
        email_out,
        password_changed_at_out,
        created_at_out;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'User not found with username: %', p_username;
    END IF;
END;
$$;
