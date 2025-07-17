CREATE TABLE api_keys (
    id INT AUTO_INCREMENT PRIMARY KEY,
    key_value VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_key_value (key_value),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

-- Insert default API keys
INSERT INTO api_keys (key_value, name) VALUES 
('test-api-key', 'Development Key'),
('production-key', 'Production Key');
