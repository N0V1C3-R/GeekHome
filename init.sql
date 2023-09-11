CREATE TABLE IF NOT EXISTS `user_entity` (
   `id` bigint NOT NULL,
   `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
   `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
   `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
   `role` tinyint(1) NOT NULL DEFAULT '1' COMMENT 'SUPER_ADMIN;ADMIN;VIP;USER',
   `active` tinyint(1) NOT NULL DEFAULT '0',
   `last_login_at` bigint NOT NULL DEFAULT '0',
   `created_at` bigint NOT NULL DEFAULT '0',
   `updated_at` bigint NOT NULL DEFAULT '0',
   `deleted_at` bigint NOT NULL DEFAULT '0',
   PRIMARY KEY (`id`),
   UNIQUE KEY `idx_username_email_status` (`username`,`email`,`deleted_at`) USING BTREE,
   UNIQUE KEY `idx_email_status` (`email`,`deleted_at`) USING BTREE,
   UNIQUE KEY `idx_username_status` (`username`,`deleted_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `user_api_key` (
    `id` bigint NOT NULL,
    `user_id` bigint NOT NULL,
    `server_name` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `api_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `is_enabled` tinyint NOT NULL DEFAULT '1',
    `created_at` bigint NOT NULL DEFAULT '0',
    `updated_at` bigint NOT NULL DEFAULT '0',
    `deleted_at` bigint NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `idx_user_server_status` (`user_id`,`server_name`,`is_enabled`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `exchange_rate` (
    `id` bigint NOT NULL,
    `trade_date` date NOT NULL,
    `foreign_currency` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `value` decimal(9,5) DEFAULT NULL,
    `is_direct_quotation` tinyint NOT NULL DEFAULT '0',
    `created_at` bigint NOT NULL DEFAULT '0',
    `updated_at` bigint NOT NULL DEFAULT '0',
    `deleted_at` bigint NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `idx_trade_date` (`trade_date`,`foreign_currency`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `blogs` (
    `id` bigint NOT NULL,
    `user_id` bigint NOT NULL,
    `title` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
    `is_anonymous` tinyint(1) NOT NULL DEFAULT '0',
    `stars` int NOT NULL DEFAULT '0',
    `total_reviews` int NOT NULL DEFAULT '0',
    `classification` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `created_at` bigint NOT NULL DEFAULT '0',
    `updated_at` bigint NOT NULL DEFAULT '0',
    `deleted_at` bigint NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `idx_title` (`title`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_stars` (`stars`) USING BTREE,
    KEY `idx_classification` (`classification`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `favorites_folder` (
    `id` bigint NOT NULL,
    `user_id` bigint NOT NULL,
    `is_folder` tinyint(1) NOT NULL,
    `nickname` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `upper` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '/',
    `path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` bigint NOT NULL DEFAULT '0',
    `updated_at` bigint NOT NULL DEFAULT '0',
    `deleted_at` bigint NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_userId_isFolder_upper_nickname_path` (`user_id`,`is_folder`,`upper`,`nickname`,`path`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
