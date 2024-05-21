DROP FUNCTION IF EXISTS create_verify_email;
CREATE OR REPLACE FUNCTION create_verify_email(
    p_user_id bigint,
    p_email VARCHAR,
    p_secret_code VARCHAR
)
RETURNS TABLE (
    id bigint,
    user_id bigint,
    email VARCHAR,
    secret_code VARCHAR,
    is_used BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE,
    expires_at_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    INSERT INTO verify_emails (user_id, email, secret_code)
    VALUES (p_user_id, p_email, p_secret_code)
    RETURNING
        verify_emails.id,
        verify_emails.user_id,
        verify_emails.email,
        verify_emails.secret_code,
        verify_emails.is_used,
        verify_emails.created_at,
        verify_emails.expired_at;
END;
$$ LANGUAGE plpgsql;
