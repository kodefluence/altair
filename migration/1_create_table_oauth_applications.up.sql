CREATE TABLE `oauth_applications` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,

  `owner_id` int(11) unsigned DEFAULT NUll,
  `description` text DEFAULT '',

  `scopes` text DEFAULT,

  `client_uid` varchar(255) NOT NULL,
  `client_secret` varchar(255) NOT NULL,

  `revoked_at` DATETIME DEFAULT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `client_uid` (`client_uid`),
  UNIQUE KEY `client_secret` (`client_secret`),
  KEY `uid_secret` (`client_uid`, `client_secret`),
  KEY `uid_secret_revoked_at` (`client_uid`, `client_secret`, `revoked_at`),
  KEY `owner_id` (`owner_id`)
) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_general_ci