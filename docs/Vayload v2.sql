CREATE TABLE `projects` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(125),
  `slug` varchar(125) UNIQUE,
  `owner_id` bigint,
  `settings` jsonb,
  `locale` varchar(10),
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `project_members` (
  `project_id` bigint,
  `user_id` bigint,
  `role_id` bigint,
  `created_at` datetime,
  PRIMARY KEY (`project_id`, `user_id`, `role_id`)
);

CREATE TABLE `users` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(125) UNIQUE,
  `email` varchar(255) UNIQUE,
  `phone` varchar(50) UNIQUE,
  `password_hash` varchar(255),
  `first_name` varchar(125),
  `last_name` varchar(125),
  `avatar_url` varchar(255),
  `email_confirmed_at` datetime,
  `phone_confirmed_at` datetime,
  `confirmed_at` datetime,
  `confirmation_token` varchar(255),
  `recovery_token` varchar(255),
  `email_change_token` varchar(255),
  `phone_change_token` varchar(255),
  `otp_code` varchar(200),
  `confirmation_sent_at` datetime,
  `recovery_sent_at` datetime,
  `email_change_sent_at` datetime,
  `phone_change_sent_at` datetime,
  `otp_sent_at` datetime,
  `email_change` varchar(255),
  `phone_change` varchar(50),
  `banned_until` datetime,
  `last_sign_in_at` datetime,
  `metadata` jsonb,
  `settings` jsonb,
  `attributes` jsonb,
  `is_super_admin` boolean,
  `is_sso_user` boolean,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime
);

CREATE TABLE `user_identities` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `user_id` bigint,
  `provider` varchar(50),
  `provider_user_id` varchar(255),
  `data` jsonb,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `sessions` (
  `id` uuid PRIMARY KEY,
  `user_id` bigint,
  `project_id` bigint,
  `ip_address` inet,
  `user_agent` text,
  `last_seen_at` datetime,
  `expires_at` datetime,
  `revoked_at` datetime,
  `created_at` datetime
);

CREATE TABLE `refresh_tokens` (
  `id` text PRIMARY KEY,
  `token_hash` text UNIQUE NOT NULL,
  `user_id` text NOT NULL,
  `family_id` text NOT NULL,
  `session_id` text,
  `parent_id` text,
  `used_at` integer,
  `revoked_at` integer,
  `revoked_reason` text,
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (unixepoch())
);

CREATE TABLE `roles` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(100) UNIQUE,
  `description` text,
  `is_system` boolean,
  `created_at` datetime
);

CREATE TABLE `permissions` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `action` varchar(50),
  `resource` varchar(100),
  `created_at` datetime
);

CREATE TABLE `role_permissions` (
  `role_id` bigint,
  `permission_id` bigint,
  PRIMARY KEY (`role_id`, `permission_id`)
);

CREATE TABLE `user_roles` (
  `user_id` bigint,
  `role_id` bigint,
  `project_id` bigint,
  PRIMARY KEY (`user_id`, `role_id`, `project_id`)
);

CREATE TABLE `collections` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `project_id` bigint,
  `name` varchar(125),
  `slug` varchar(125),
  `description` text,
  `icon` varchar(50),
  `is_system` boolean,
  `supports_versioning` boolean,
  `supports_localization` boolean,
  `settings` jsonb,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime
);

CREATE TABLE `collection_fields` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `collection_id` bigint,
  `name` varchar(125),
  `slug` varchar(125),
  `field_type` varchar(50),
  `is_required` boolean,
  `is_unique` boolean,
  `is_localized` boolean,
  `is_indexed` boolean,
  `default_value` text,
  `validation_rules` jsonb,
  `options` jsonb,
  `display_order` int,
  `help_text` text,
  `placeholder` text,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `entries` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `project_id` bigint,
  `collection_id` bigint,
  `author_id` bigint,
  `status` varchar(20),
  `version` int,
  `parent_version_id` bigint,
  `locale` varchar(10),
  `title` varchar(500),
  `slug` varchar(255),
  `published_at` datetime,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime
);

CREATE TABLE `entry_fields` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `entry_id` bigint,
  `field_id` bigint,
  `locale` varchar(10),
  `value` text,
  `value_json` jsonb,
  `value_number` numeric,
  `value_boolean` boolean,
  `value_date` datetime,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `entry_relations` (
  `source_entry_id` bigint,
  `target_entry_id` bigint,
  `field_id` bigint,
  `relation_order` int,
  PRIMARY KEY (`source_entry_id`, `target_entry_id`, `field_id`)
);

CREATE TABLE `folder_objects` (
  `id` text PRIMARY KEY,
  `owner_id` text NOT NULL,
  `project_id` text NOT NULL,
  `parent_id` text,
  `name` text NOT NULL,
  `path` text,
  `depth` int,
  `file_count` int DEFAULT 0,
  `subfolder_count` int DEFAULT 0,
  `created_at` timestamp,
  `updated_at` timestamp
);

CREATE TABLE `file_objects` (
  `id` text PRIMARY KEY,
  `owner_id` text NOT NULL,
  `project_id` text NOT NULL,
  `folder_id` text,
  `name` text NOT NULL,
  `mime_type` text,
  `category` text,
  `size` bigint,
  `provider` text,
  `provider_key` text,
  `metadata` json,
  `created_at` timestamp,
  `updated_at` timestamp
);

CREATE TABLE `file_usage` (
  `file_id` bigint,
  `entry_id` bigint,
  `field_id` bigint,
  `usage_order` int,
  PRIMARY KEY (`file_id`, `entry_id`, `field_id`)
);

CREATE TABLE `plugins` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `name` varchar(255),
  `version` varchar(50),
  `description` text,
  `author_name` varchar(255),
  `author_email` varchar(255),
  `author_homepage` varchar(500),
  `entry_point` varchar(500),
  `dependencies` jsonb,
  `manifest` jsonb,
  `checksum` text,
  `is_active` boolean,
  `is_system` boolean,
  `created_at` datetime,
  `updated_at` datetime
);

CREATE TABLE `plugin_hooks_log` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `plugin_id` bigint,
  `hook_name` varchar(100),
  `execution_time_ms` int,
  `success` boolean,
  `error_message` text,
  `created_at` datetime
);

CREATE TABLE `notifications` (
  `id` bigint PRIMARY KEY,
  `title` varchar(255),
  `body` text,
  `datetime` datetime,
  `status` enum(sent,read,unread,dismissed,failed),
  `type` varchar(255),
  `actions` json,
  `sent_at` datetime,
  `read_at` datetime,
  `user_id` bigint,
  `project_id` bigint,
  `created_at` datetime,
  `updated_at` datetime,
  `deleted_at` datetime
);

CREATE TABLE `notification_logs` (
  `id` bigint,
  `type` varchar(100),
  `sender` varchar(255),
  `channel` varchar(20),
  `destination` varchar(255),
  `url` varchar(255),
  `request_body` text,
  `response_body` text,
  `response_status` int,
  `error` text,
  `notifications_id` varchar(50),
  `created_at` datetime
);

CREATE TABLE `activity_logs` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `project_id` bigint,
  `actor_id` bigint,
  `action` varchar(100),
  `resource_type` varchar(100),
  `resource_id` bigint,
  `resource_slug` varchar(255),
  `message` text,
  `changes` jsonb,
  `payload` jsonb,
  `origin` varchar(50),
  `severity` varchar(20),
  `request_id` uuid,
  `ip_address` inet,
  `user_agent` text,
  `created_at` datetime
);

CREATE UNIQUE INDEX `permissions_index_0` ON `permissions` (`action`, `resource`);

CREATE UNIQUE INDEX `collections_index_1` ON `collections` (`project_id`, `slug`);

CREATE UNIQUE INDEX `collection_fields_index_2` ON `collection_fields` (`collection_id`, `slug`);

CREATE INDEX `entries_index_3` ON `entries` (`collection_id`, `slug`, `locale`);

CREATE UNIQUE INDEX `entry_fields_index_4` ON `entry_fields` (`entry_id`, `field_id`, `locale`);

CREATE INDEX `entry_fields_index_5` ON `entry_fields` (`field_id`);

CREATE INDEX `activity_logs_index_6` ON `activity_logs` (`project_id`, `created_at`);

CREATE INDEX `activity_logs_index_7` ON `activity_logs` (`actor_id`, `created_at`);

CREATE INDEX `activity_logs_index_8` ON `activity_logs` (`resource_type`, `resource_id`);

CREATE INDEX `activity_logs_index_9` ON `activity_logs` (`request_id`);

ALTER TABLE `refresh_tokens` ADD FOREIGN KEY (`parent_id`) REFERENCES `refresh_tokens` (`id`) ON DELETE SET NULL;

ALTER TABLE `folder_objects` ADD FOREIGN KEY (`parent_id`) REFERENCES `folder_objects` (`id`);

ALTER TABLE `file_objects` ADD FOREIGN KEY (`folder_id`) REFERENCES `folder_objects` (`id`);

ALTER TABLE `notifications` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `notifications` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `collections` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `collection_fields` ADD FOREIGN KEY (`collection_id`) REFERENCES `collections` (`id`);

ALTER TABLE `entries` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `entries` ADD FOREIGN KEY (`collection_id`) REFERENCES `collections` (`id`);

ALTER TABLE `entries` ADD FOREIGN KEY (`author_id`) REFERENCES `users` (`id`);

ALTER TABLE `entries` ADD FOREIGN KEY (`parent_version_id`) REFERENCES `entries` (`id`);

ALTER TABLE `entry_fields` ADD FOREIGN KEY (`entry_id`) REFERENCES `entries` (`id`);

ALTER TABLE `entry_fields` ADD FOREIGN KEY (`field_id`) REFERENCES `collection_fields` (`id`);

ALTER TABLE `entry_relations` ADD FOREIGN KEY (`source_entry_id`) REFERENCES `entries` (`id`);

ALTER TABLE `entry_relations` ADD FOREIGN KEY (`target_entry_id`) REFERENCES `entries` (`id`);

ALTER TABLE `entry_relations` ADD FOREIGN KEY (`field_id`) REFERENCES `collection_fields` (`id`);

ALTER TABLE `file_usage` ADD FOREIGN KEY (`file_id`) REFERENCES `file_objects` (`id`);

ALTER TABLE `file_usage` ADD FOREIGN KEY (`entry_id`) REFERENCES `entries` (`id`);

ALTER TABLE `file_usage` ADD FOREIGN KEY (`field_id`) REFERENCES `collection_fields` (`id`);

ALTER TABLE `projects` ADD FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`);

ALTER TABLE `project_members` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `project_members` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `project_members` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

ALTER TABLE `user_identities` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `sessions` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `refresh_tokens` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `refresh_tokens` ADD FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`);

ALTER TABLE `role_permissions` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

ALTER TABLE `role_permissions` ADD FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `file_objects` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `file_objects` ADD FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`);

ALTER TABLE `plugin_hooks_log` ADD FOREIGN KEY (`plugin_id`) REFERENCES `plugins` (`id`);
