CREATE TABLE `oauth_access_tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,

  `oauth_application_id` int(11) unsigned NOT NULL,
  `resource_owner_id` int(11) unsigned NOT NULL,

  `token` varchar(255) NOT NULL,
  `scopes` text,

  `expires_in` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL,
  `revoked_at` DATETIME DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`),
  KEY `id_oauth_application_id` (`id`, `oauth_application_id`),
  KEY `id_oauth_application_id_resource_owner_id` (`id`, `oauth_application_id`, `resource_owner_id`),
  KEY `oauth_application_id_resource_owner_id` (`oauth_application_id`, `resource_owner_id`)

) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_general_ci;