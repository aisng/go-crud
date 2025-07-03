CREATE TABLE `users` (
    `id` bigint(20) AUTO_INCREMENT PRIMARY KEY,
    `email` varchar(100) NOT NULL UNIQUE,
    `username` varchar(100) NOT NULL UNIQUE,
    `password` varchar(72) NOT NULL,
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
