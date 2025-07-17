ALTER TABLE users 
DROP COLUMN reset_password_expires_at,
CHANGE COLUMN reset_password_token verification_token VARCHAR(255);
