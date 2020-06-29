CREATE TABLE `eos_auth_tokens` (
    `token` varchar(43) NOT NULL,
    `token_sha256` varchar(64) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `eos_account` char(12) NOT NULL,
    `status` tinyint(4) NOT NULL COMMENT '0 初始 1 已绑定账号 2 通知成功 3 通知失败',
    `notice_error` varchar(1000) NOT NULL,
    PRIMARY KEY (`token`),
    UNIQUE KEY `token_sha256_UNIQUE` (`token_sha256`),
    KEY `IX_eosAccount` (`eos_account`),
    KEY `IX_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `eos_deposits` (
    `id` int(11) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `amount` decimal(10,4) NOT NULL,
    `from` varchar(12) NOT NULL,
    `to` varchar(12) NOT NULL,
    `memo` varchar(256) NOT NULL,
    `transaction_id` char(64) NOT NULL,
    `block_confirmed` tinyint(4) NOT NULL,
    `status` tinyint(4) NOT NULL COMMENT '0 初始 1 通知成功 2 通知失败',
    `transaction_ms` decimal(13,0) NOT NULL,
    `notice_error` varchar(1000) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `IX_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE `eos_withdraws` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `amount` decimal(10,4) NOT NULL,
    `from` varchar(45) NOT NULL,
    `to` varchar(45) NOT NULL,
    `memo` varchar(45) NOT NULL,
    `transaction_id` char(64) NOT NULL,
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 初始 1 已发送未确认 2 已确认',
    `consumer_sid` varchar(45) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `UK_consumerSid` (`consumer_sid`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

