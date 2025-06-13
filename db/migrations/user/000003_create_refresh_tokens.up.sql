CREATE TABLE refresh_tokens (
                                id VARCHAR(36) PRIMARY KEY,
                                user_id VARCHAR(36) NOT NULL,
                                token VARCHAR(255) NOT NULL,
                                user_agent VARCHAR(255),
                                ip_address VARCHAR(45),
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                expires_at TIMESTAMP NOT NULL,
                                revoked TINYINT(1) DEFAULT 0,
                                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                INDEX idx_user_id (user_id),
                                INDEX idx_token (token)
);