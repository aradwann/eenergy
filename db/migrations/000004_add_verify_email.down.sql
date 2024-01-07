-- Drop the foreign key constraint on 'verify_emails'
ALTER TABLE verify_emails
DROP CONSTRAINT IF EXISTS fk_verify_emails_users;

-- Drop the 'verify_emails' table
DROP TABLE IF EXISTS verify_emails;

-- Remove the 'is_email_verified' column from 'users'
ALTER TABLE users
DROP COLUMN IF EXISTS is_email_verified;
