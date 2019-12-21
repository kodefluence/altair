CREATE TABLE `oauth_access_tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,

  `oauth_application_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL,

  `token` varchar(255) NOT NULL,

  `expires_in` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL,
  `revoked_at` DATETIME DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`),
  KEY `id_oauth_application_id` (`id`, `oauth_application_id`),
  KEY `id_oauth_application_id_user_id` (`id`, `oauth_application_id`, `user_id`),
  KEY `oauth_application_id_user_id` (`oauth_application_id`, `user_id`)

) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_general_ci;