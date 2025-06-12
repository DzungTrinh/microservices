CREATE TABLE IF NOT EXISTS users (
                                     id CHAR(36) PRIMARY KEY, -- CHAR(36) as string
                                     created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     username VARCHAR(50) UNIQUE NOT NULL,
                                     email VARCHAR(100) UNIQUE NOT NULL,
                                     password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS roles (
                                     id CHAR(36) PRIMARY KEY,
                                     name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_roles (
                                          user_id CHAR(36) NOT NULL,
                                          role_id CHAR(36) NOT NULL,
                                          PRIMARY KEY (user_id, role_id),
                                          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                          FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);