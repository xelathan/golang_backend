CREATE TABLE IF NOT EXISTS `user_addresses` (
    `userId` INT UNSIGNED NOT NULL,
    `default` TEXT NOT NULL,
    `secondary` TEXT NOT NULL,
    `tertiary` TEXT NOT NULL,

    PRIMARY KEY (`userId`),
    FOREIGN KEY (`userId`) REFERENCES users(`id`)
);