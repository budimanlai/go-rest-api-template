CREATE TABLE `api_key` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `description` text DEFAULT NULL,
  `api_key` varchar(32) NOT NULL,
  `auth_key` varchar(64) NOT NULL,
  `status` varchar(15) NOT NULL DEFAULT 'active',
  `h2h` char(1) NOT NULL DEFAULT 'N',
  `last_access` datetime DEFAULT NULL,
  `ip_whitelist` varchar(256) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `created_by` int(11) unsigned NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_api_key` (`api_key`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
