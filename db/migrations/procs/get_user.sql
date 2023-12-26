DROP PROCEDURE IF EXISTS get_user;
CREATE OR REPLACE PROCEDURE get_user(
    IN p_username VARCHAR,
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
    SELECT
        users.username,
        users.hashed_password,
        users.fullname,
        users.email,
        users.password_changed_at,
        users.created_at
    INTO
        username_out,
        hashed_password_out,
        fullname_out,
        email_out,
        password_changed_at_out,
        created_at_out
    FROM users
    WHERE users.username = p_username
    LIMIT 1;
END;
$$;
