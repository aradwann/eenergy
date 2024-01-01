DROP FUNCTION IF EXISTS get_session;
CREATE OR REPLACE FUNCTION get_session(
    p_id UUID
)
RETURNS TABLE (
    id UUID,
    username VARCHAR,
    refresh_token VARCHAR,
    user_agent VARCHAR,
    client_ip VARCHAR,
    is_blocked BOOLEAN,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        sessions.id,
        sessions.username,
        sessions.refresh_token,
        sessions.user_agent,
        sessions.client_ip,
        sessions.is_blocked,
        sessions.expires_at,
        sessions.created_at
    FROM sessions
    WHERE sessions.id = p_id
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;
