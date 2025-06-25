CREATE DATABASE user_service;
GRANT ALL PRIVILEGES ON user_service.* TO 'user'@'%';

CREATE DATABASE rbac_service;
GRANT ALL PRIVILEGES ON rbac_service.* TO 'user'@'%';

CREATE DATABASE notification_service;
GRANT ALL PRIVILEGES ON notification_service.* TO 'user'@'%';

FLUSH PRIVILEGES;