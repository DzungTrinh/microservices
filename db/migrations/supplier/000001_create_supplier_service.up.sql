CREATE TABLE suppliers (
                         id CHAR(36) PRIMARY KEY,
                         name VARCHAR(100) NOT NULL,
                         description TEXT,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menu_items (
                            id CHAR(36) PRIMARY KEY,
                            supplier_id CHAR(36) REFERENCES suppliers(id) ON DELETE CASCADE,
                            name VARCHAR(100) NOT NULL,
                            description TEXT,
                            price DECIMAL(10, 2) NOT NULL,
                            available BOOLEAN DEFAULT TRUE,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);