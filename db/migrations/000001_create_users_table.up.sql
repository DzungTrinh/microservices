CREATE TABLE users (
                       id BIGINT AUTO_INCREMENT PRIMARY KEY,
                       created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL
);