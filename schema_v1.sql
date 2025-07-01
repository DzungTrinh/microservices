-- ─────────────────────────────────────────
-- identity-svc schema (identity_db)
-- ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users
(
    id             CHAR(36) PRIMARY KEY,
    email          VARCHAR(100) UNIQUE NOT NULL,
    username       VARCHAR(50) UNIQUE  NOT NULL,
    email_verified TINYINT(1)          NOT NULL DEFAULT 0,
    created_at     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at     DATETIME            NULL
);

CREATE TABLE IF NOT EXISTS credentials
(
    id           CHAR(36) PRIMARY KEY,
    user_id      CHAR(36)    NOT NULL,
    provider     VARCHAR(20) NOT NULL,
    secret_hash  VARCHAR(255),
    provider_uid VARCHAR(255),
    created_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   DATETIME    NULL,
    UNIQUE KEY uniq_provider_uid (provider, provider_uid),
    CONSTRAINT fk_cred_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mfa_factors
(
    id         CHAR(36) PRIMARY KEY,
    user_id    CHAR(36)    NOT NULL,
    type       VARCHAR(20) NOT NULL,
    secret     TEXT,
    verified   TINYINT(1)  NOT NULL DEFAULT 0,
    created_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME    NULL,
    CONSTRAINT fk_mfa_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS verification_tokens
(
    token      CHAR(64) PRIMARY KEY,
    user_id    CHAR(36)    NOT NULL,
    channel    VARCHAR(30) NOT NULL,
    expires_at DATETIME    NOT NULL,
    created_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME    NULL,
    INDEX idx_expires (expires_at),
    CONSTRAINT fk_vtoken_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    id         CHAR(36) PRIMARY KEY,
    user_id    CHAR(36)   NOT NULL,
    token      CHAR(128)  NOT NULL,
    user_agent VARCHAR(255),
    ip_address VARCHAR(45),
    created_at DATETIME   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME   NOT NULL,
    revoked    TINYINT(1) NOT NULL DEFAULT 0,
    deleted_at DATETIME   NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_token (token),
    CONSTRAINT fk_rt_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS device_sessions
(
    id         CHAR(36) PRIMARY KEY,
    user_id    CHAR(36) NOT NULL,
    user_agent VARCHAR(255),
    ip_address VARCHAR(45),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME,
    deleted_at DATETIME NULL,
    INDEX idx_ds_user (user_id),
    CONSTRAINT fk_ds_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_profiles
(
    user_id      CHAR(36) PRIMARY KEY,
    profile      JSON     NOT NULL,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at   DATETIME NULL,
    display_name VARCHAR(191) GENERATED ALWAYS AS (LOWER(JSON_UNQUOTE(JSON_EXTRACT(profile, '$.displayName')))) VIRTUAL,
    INDEX idx_display_name (display_name),
    CONSTRAINT fk_profile_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- ─────────────────────────────────────────
-- rbac-svc schema (rbac_db)
-- ─────────────────────────────────────────
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
    role_id    CHAR(36) NOT NULL,
    perm_id    CHAR(36) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    PRIMARY KEY (role_id, perm_id),
    CONSTRAINT fk_rp_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    CONSTRAINT fk_rp_perm FOREIGN KEY (perm_id) REFERENCES permissions (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_roles
(
    user_id    CHAR(36) NOT NULL,
    role_id    CHAR(36) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
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
    deleted_at DATETIME NULL,
    PRIMARY KEY (user_id, perm_id),
    CONSTRAINT fk_up_perm FOREIGN KEY (perm_id) REFERENCES permissions (id) ON DELETE CASCADE
);

-- ─────────────────────────────────────────
-- notification-svc schema (notification_db)
-- ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS templates
(
    id         CHAR(36) PRIMARY KEY,
    channel    VARCHAR(20)         NOT NULL,
    name       VARCHAR(100) UNIQUE NOT NULL,
    subject    VARCHAR(255),
    body_html  TEXT,
    body_text  TEXT,
    version    INT                 NOT NULL DEFAULT 1,
    created_at DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME            NULL
);

CREATE TABLE IF NOT EXISTS messages
(
    id          CHAR(36) PRIMARY KEY,
    user_id     CHAR(36)                      NULL,
    template_id CHAR(36)                      NULL,
    payload     JSON                          NULL,
    status      ENUM ('queued','sent','fail') NOT NULL DEFAULT 'queued',
    retries     TINYINT                       NOT NULL DEFAULT 0,
    created_at  DATETIME                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sent_at     DATETIME                      NULL,
    deleted_at  DATETIME                      NULL,
    INDEX idx_status (status),
    CONSTRAINT fk_msg_template FOREIGN KEY (template_id) REFERENCES templates (id) ON DELETE SET NULL
);

-- ─────────────────────────────────────────
-- audit-svc schema (audit_db)
-- ─────────────────────────────────────────
CREATE TABLE IF NOT EXISTS audit_events
(
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    occurred_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actor_id    CHAR(36)     NULL,
    action      VARCHAR(100) NOT NULL,
    resource    VARCHAR(255) NULL,
    ip_address  VARCHAR(45)  NULL,
    meta        JSON         NULL,
    INDEX idx_actor (actor_id),
    INDEX idx_action (action),
    INDEX idx_resource (resource),
    INDEX idx_occurred (occurred_at)
) ENGINE = InnoDB
    PARTITION BY RANGE (TO_DAYS(occurred_at)) (
        PARTITION p2025_06 VALUES LESS THAN (TO_DAYS('2025-07-01')),
        PARTITION pmax VALUES LESS THAN MAXVALUE
        );
