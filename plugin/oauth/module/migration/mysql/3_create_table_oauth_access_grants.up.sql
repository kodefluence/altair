CREATE TABLE `oauth_access_grants` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,

  `oauth_application_id` int(11) unsigned NOT NULL,
  `resource_owner_id` int(11) unsigned NOT NULL,
  `scopes` text,

  `code` varchar(255) NOT NULL,
  `redirect_uri` text,

  `expires_in` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL,
  `revoked_at` DATETIME DEFAULT NULL,

  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  KEY `id_oauth_application_id` (`id`, `oauth_application_id`),
  KEY `id_oauth_application_id_resource_owner_id` (`id`, `oauth_application_id`, `resource_owner_id`),
  KEY `oauth_application_id_resource_owner_id` (`oauth_application_id`, `resource_owner_id`)

) ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_general_ci;