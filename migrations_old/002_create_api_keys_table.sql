-- Migration: Create API keys table
-- Created at: 2025-07-17 10:05:00

-- +migrate Up
-- Create API keys table for authentication
CREATE TABLE IF NOT EXISTS api_keys (
    id INT AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status ENUM('active', 'inactive') DEFAULT 'active',
    expires_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NULL,
    
    INDEX idx_key (`key`),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_created_at (created_at),
    INDEX idx_last_used (last_used_at)
);

-- Insert default API key for development
INSERT INTO api_keys (`key`, name, description, status) VALUES 
('test-api-key', 'Development Key', 'Default API key for development and testing', 'active');

-- +migrate Down
-- Drop API keys table
DROP TABLE IF EXISTS api_keys;
