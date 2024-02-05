DROP FUNCTION IF EXISTS create_verify_email;
CREATE OR REPLACE FUNCTION create_verify_email(
    p_username VARCHAR,
    p_email VARCHAR,
    p_secret_code VARCHAR
)
RETURNS TABLE (
    id bigint,
    username VARCHAR,
    email VARCHAR,
    secret_code VARCHAR,
    is_used BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE,
    expires_at_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    INSERT INTO verify_emails (username, email, secret_code)
    VALUES (p_username, p_email, p_secret_code)
    RETURNING
        verify_emails.id,
        verify_emails.username,
        verify_emails.email,
        verify_emails.secret_code,
        verify_emails.is_used,
        verify_emails.created_at,
        verify_emails.expired_at;
END;
$$ LANGUAGE plpgsql;
