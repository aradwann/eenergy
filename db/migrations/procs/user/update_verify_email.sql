DROP FUNCTION IF EXISTS update_verify_email;
CREATE OR REPLACE FUNCTION update_verify_email(
    p_id bigint,
    p_secret_code VARCHAR
)
RETURNS TABLE (
    id bigint,
    username VARCHAR,
    email VARCHAR,
    secret_code VARCHAR,
    is_used BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE,
    expires_at_at TIMESTAMP WITH TIME ZONE,
    role VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    UPDATE verify_emails
    SET is_used = TRUE
    WHERE 
    verify_emails.id = p_id
    AND verify_emails.secret_code = p_secret_code
    AND verify_emails.is_used = FALSE
    AND verify_emails.expired_at > now()
    RETURNING
        verify_emails.id,
        verify_emails.username,
        verify_emails.email,
        verify_emails.secret_code,
        verify_emails.is_used,
        verify_emails.created_at,
        verify_emails.expired_at,
        verify_emails.role;

END;
$$ LANGUAGE plpgsql;
