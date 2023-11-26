CREATE OR REPLACE FUNCTION CreateUser(
    p_username VARCHAR,
    p_hashed_password VARCHAR,
    p_fullname VARCHAR,
    p_email VARCHAR
)
RETURNS TABLE (
    username VARCHAR,
    hashed_password VARCHAR,
    fullname VARCHAR,
    email VARCHAR
) AS $$
BEGIN
    INSERT INTO users (username, hashed_password, fullname, email)
    VALUES (p_username, p_hashed_password, p_fullname, p_email)
    RETURNING *;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION GetUser(
    p_username VARCHAR
)
RETURNS TABLE (
    username VARCHAR,
    hashed_password VARCHAR,
    fullname VARCHAR,
    email VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM users
    WHERE username = p_username
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION UpdateUser(
    p_username VARCHAR,
    p_hashed_password VARCHAR,
    p_password_changed_at TIMESTAMP,
    p_fullname VARCHAR,
    p_email VARCHAR
)
RETURNS TABLE (
    username VARCHAR,
    hashed_password VARCHAR,
    fullname VARCHAR,
    email VARCHAR
) AS $$
BEGIN
    UPDATE users
    SET
        hashed_password = COALESCE(p_hashed_password, hashed_password),
        password_changed_at = COALESCE(p_password_changed_at, password_changed_at),
        fullname = COALESCE(p_fullname, fullname),
        email = COALESCE(p_email, email)
    WHERE username = p_username
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
