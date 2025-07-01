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
