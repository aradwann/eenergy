DROP FUNCTION IF EXISTS delete_account;
CREATE OR REPLACE FUNCTION delete_account(p_id INT, OUT result BOOLEAN)
RETURNS BOOLEAN
AS $$
BEGIN
    DELETE FROM accounts WHERE id = p_id;
    IF FOUND THEN
        result := TRUE;
    ELSE
        result := FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;
