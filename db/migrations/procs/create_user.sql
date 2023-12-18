CREATE OR REPLACE PROCEDURE create_user(
    p_username VARCHAR,
    p_hashed_password VARCHAR,
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
    INSERT INTO users (username, hashed_password, fullname, email)
    VALUES (p_username, p_hashed_password, p_fullname, p_email)
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
END;
$$;
