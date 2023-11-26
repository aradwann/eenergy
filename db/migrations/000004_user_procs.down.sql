-- Drop CreateUser stored procedure
DROP FUNCTION IF EXISTS CreateUser(VARCHAR, VARCHAR, VARCHAR, VARCHAR);

-- Drop GetUser stored procedure
DROP FUNCTION IF EXISTS GetUser(VARCHAR);

-- Drop UpdateUser stored procedure
DROP FUNCTION IF EXISTS UpdateUser(VARCHAR, VARCHAR, TIMESTAMP, VARCHAR, VARCHAR);
