CREATE DATABASE user_service;
GRANT ALL PRIVILEGES ON user_service.* TO 'user'@'%';

CREATE DATABASE order_service;
GRANT ALL PRIVILEGES ON order_service.* TO 'user'@'%';

CREATE DATABASE payment_service;
GRANT ALL PRIVILEGES ON payment_service.* TO 'user'@'%';

CREATE DATABASE supplier_service;
GRANT ALL PRIVILEGES ON supplier_service.* TO 'user'@'%';

CREATE DATABASE notification_service;
GRANT ALL PRIVILEGES ON notification_service.* TO 'user'@'%';

FLUSH PRIVILEGES;