CREATE TABLE IF NOT EXISTS `user_addresses` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `userId` INT UNSIGNED NOT NULL,
    `address_type` ENUM('first', 'secondary', 'tertiary') NOT NULL,
    `address` TEXT NOT NULL,
    FOREIGN KEY (`userId`) REFERENCES users(`id`),
    UNIQUE KEY `unique_user_address` (`userId`, `address_type`)
);