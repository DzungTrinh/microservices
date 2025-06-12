CREATE TABLE IF NOT EXISTS notifications (
                                             id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                             user_id BIGINT NOT NULL,
                                             type VARCHAR(50) NOT NULL, -- e.g., "order_status", "promo", "message"
                                             title VARCHAR(255),
                                             message TEXT,
                                             is_read BOOLEAN DEFAULT FALSE,
                                             created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);