-- Sample API keys for testing
INSERT INTO `api_key` (`id`, `name`, `description`, `api_key`, `auth_key`, `status`, `h2h`, `last_access`, `ip_whitelist`, `created_at`, `created_by`, `updated_at`, `updated_by`) VALUES
(1, 'Development API Key', 'API key for development environment', 'dev_api_key_12345678901234567890', 'dev_auth_key_1234567890123456789012345678901234567890123456789012', 'active', 'N', NULL, NULL, NOW(), 1, NOW(), NULL),
(2, 'Production API Key', 'API key for production environment', 'prod_api_key_87654321098765432109', 'prod_auth_key_0987654321098765432109876543210987654321098765432109', 'active', 'Y', NULL, '127.0.0.1,192.168.1.1', NOW(), 1, NOW(), NULL),
(3, 'Partner API Key', 'API key for external partner integration', 'partner_api_key_11111111111111111111', 'partner_auth_key_1111111111111111111111111111111111111111111111111111', 'inactive', 'N', NULL, NULL, NOW(), 1, NOW(), NULL);
