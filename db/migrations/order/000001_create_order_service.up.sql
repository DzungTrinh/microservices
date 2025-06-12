CREATE TABLE orders (
                        id CHAR(36) PRIMARY KEY,
                        user_id CHAR(36) NOT NULL, -- from user-service
                        vendor_id CHAR(36) NOT NULL, -- from vendor-service
                        status VARCHAR(20) NOT NULL, -- pending, preparing, ready, completed, cancelled
                        total_price DECIMAL(10, 2) NOT NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
                             id CHAR(36) PRIMARY KEY,
                             order_id CHAR(36) REFERENCES orders(id) ON DELETE CASCADE,
                             menu_item_id CHAR(36) NOT NULL,
                             quantity INTEGER NOT NULL,
                             price DECIMAL(10, 2) NOT NULL
);
