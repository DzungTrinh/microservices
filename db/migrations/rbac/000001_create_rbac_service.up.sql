CREATE TABLE IF NOT EXISTS roles
(
    id         CHAR(36) PRIMARY KEY,
    name       VARCHAR(50) UNIQUE NOT NULL,
    built_in   TINYINT(1)         NOT NULL DEFAULT 0,
    created_at DATETIME           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME           NULL
);

CREATE TABLE IF NOT EXISTS permissions
(
    id         CHAR(36) PRIMARY KEY,
    name       VARCHAR(100) UNIQUE NOT NULL,
    created_at DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME            NULL
);

CREATE TABLE IF NOT EXISTS role_permissions
(
    role_id CHAR(36) NOT NULL,
    perm_id CHAR(36) NOT NULL,
    PRIMARY KEY (role_id, perm_id),
    CONSTRAINT fk_rp_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    CONSTRAINT fk_rp_perm FOREIGN KEY (perm_id) REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_roles
(
    user_id CHAR(36) NOT NULL,
    role_id CHAR(36) NOT NULL,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_ur_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_permissions
(
    user_id    CHAR(36) NOT NULL,
    perm_id    CHAR(36) NOT NULL,
    granter_id CHAR(36) NULL,
    expires_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, perm_id),
    CONSTRAINT fk_up_perm FOREIGN KEY (perm_id) REFERENCES permissions (id) ON DELETE CASCADE
);
