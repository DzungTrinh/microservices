CREATE TABLE payments (
                          id CHAR(36) PRIMARY KEY,
                          order_id CHAR(36) NOT NULL,
                          amount DECIMAL(10, 2) NOT NULL,
                          payment_method VARCHAR(50),
                          status VARCHAR(20) NOT NULL, -- pending, paid, failed, refunded
                          paid_at TIMESTAMP
);