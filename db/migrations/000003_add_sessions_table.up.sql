CREATE TABLE `sessions` (
    `id` VARCHAR(50) PRIMARY KEY NOT NULL,
    `user_email` VARCHAR(50) NOT NULL,
    `refresh_token` VARCHAR(512) NOT NULL,
    `is_revoked` BOOL NOT NULL DEFAULT false,
    `created_at` DATETIME DEFAULT now(),
    `expires_at` DATETIME 
);