-- Migration: Create users table
-- Created at: 2025-07-17 10:00:00

-- +migrate Up
-- Create users table with all necessary fields
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    reset_password_token VARCHAR(255) NULL,
    reset_password_expires_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_by INT NULL,
    updated_by INT NULL,
    deleted_by INT NULL,
    
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_status (status),
    INDEX idx_reset_token (reset_password_token),
    INDEX idx_reset_expires (reset_password_expires_at),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at)
);

-- +migrate Down
-- Drop users table
DROP TABLE IF EXISTS users;
