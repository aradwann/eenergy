DROP FUNCTION IF EXISTS create_session;
CREATE OR REPLACE FUNCTION create_session(
    p_id UUID,
    p_user_id bigint,
    p_refresh_token VARCHAR,
    p_user_agent VARCHAR,
    p_client_ip VARCHAR,
    p_is_blocked BOOLEAN,
    p_expires_at TIMESTAMP WITH TIME ZONE,
    p_created_at TIMESTAMP WITH TIME ZONE
)
RETURNS TABLE (
    id UUID,
    user_id bigint,
    refresh_token VARCHAR,
    user_agent VARCHAR,
    client_ip VARCHAR,
    is_blocked BOOLEAN,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    -- Insert a new session for the provided parameters
    RETURN QUERY
    INSERT INTO sessions (
       id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
    )
    VALUES (
       p_id, p_user_id, p_refresh_token, p_user_agent, p_client_ip, p_is_blocked, p_expires_at, p_created_at
    )
    RETURNING *;
END;
$$ LANGUAGE plpgsql;
