-- +goose Up
-- SQL in this section is executed when the migration is applied
CREATE TABLE IF NOT EXISTS `users` (
  `id` VARCHAR(32) NOT NULL PRIMARY KEY,
  `username` VARCHAR(50) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(100) NOT NULL UNIQUE,
  `full_name` VARCHAR(100) NOT NULL,
  `organization` VARCHAR(100) NOT NULL,
  `role` VARCHAR(20) NOT NULL DEFAULT 'student',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `last_login_at` TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建索引
CREATE INDEX `idx_users_role` ON `users` (`role`);
CREATE INDEX `idx_users_created_at` ON `users` (`created_at`);

-- 创建管理员用户 (密码: admin123)
INSERT INTO `users` (`id`, `username`, `password`, `email`, `full_name`, `organization`, `role`)
VALUES ('u00000001', 'admin', '$2a$10$aZB36UooZpL.fAgbQVN/j.pfZVVvkHxEnj7vfkVSqwBr/setGfNgu', 'admin@example.com', '系统管理员', '海洋研究所', 'admin');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP TABLE IF EXISTS `users`; 