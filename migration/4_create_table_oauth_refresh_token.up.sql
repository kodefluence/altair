CREATE TABLE `oauth_refresh_tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,

  `oauth_access_token_id` int(11) unsigned NOT NULL,

  `token` varchar(255) NOT NULL,

  `expires_in` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL,
  `revoked_at` DATETIME DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`),
  KEY `oauth_access_token_id` (`id`, `oauth_access_token_id`)

) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_general_ci;
