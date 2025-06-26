# CREATE TABLE IF NOT EXISTS users (
#      id CHAR(36) PRIMARY KEY, -- CHAR(36) as string
#      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
#      updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
#      username VARCHAR(50) UNIQUE NOT NULL,
#      email VARCHAR(100) UNIQUE NOT NULL,
#      password VARCHAR(255) NOT NULL
# );
#
# CREATE TABLE IF NOT EXISTS roles (
#      id CHAR(36) PRIMARY KEY,
#      name VARCHAR(50) UNIQUE NOT NULL
# );
#
# CREATE TABLE IF NOT EXISTS user_roles (
#       user_id CHAR(36) NOT NULL,
#       role_id CHAR(36) NOT NULL,
#       PRIMARY KEY (user_id, role_id),
#       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
#       FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
# );
#
# INSERT INTO roles (id, name) VALUES (UUID(), 'user'), (UUID(), 'admin');
#
# CREATE TABLE refresh_tokens (
#     id VARCHAR(36) PRIMARY KEY,
#     user_id VARCHAR(36) NOT NULL,
#     token TEXT NOT NULL,
#     user_agent VARCHAR(255),
#     ip_address VARCHAR(45),
#     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
#     expires_at TIMESTAMP NOT NULL,
#     revoked TINYINT(1) NOT NULL DEFAULT 0,
#     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
#     INDEX idx_user_id (user_id),
#     INDEX idx_token (token(255))
# );

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

CREATE TABLE IF NOT EXISTS outbox_events
(
    id             BIGINT AUTO_INCREMENT PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id   VARCHAR(100) NOT NULL,
    type           VARCHAR(100) NOT NULL,
    payload        JSON         NOT NULL,
    status         VARCHAR(20)  NOT NULL DEFAULT 'pending',
    created_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at   DATETIME
);

CREATE INDEX idx_outbox_status ON outbox_events (status);