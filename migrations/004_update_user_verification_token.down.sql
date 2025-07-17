ALTER TABLE users 
ADD COLUMN reset_password_expires_at TIMESTAMP NULL,
CHANGE COLUMN verification_token reset_password_token VARCHAR(255);
