package constants

import "time"

const (
	GrpcPort       = "GRPC_PORT"
	HttpPort       = "HTTP_PORT"
	ConfigPath     = "CONFIG_PATH"
	KafkaBrokers   = "KAFKA_BROKERS"
	JaegerHostPort = "JAEGER_HOST"
	RedisAddr      = "REDIS_ADDR"
	MongoDbURI     = "MONGO_URI"
	PostgresqlHost = "POSTGRES_HOST"
	PostgresqlPort = "POSTGRES_PORT"
	MysqlHost      = "MYSQL_HOST"
	MysqlPort      = "MYSQL_PORT"

	ReaderServicePort = "READER_SERVICE"

	Yaml     = "yaml"
	Redis    = "redis"
	Kafka    = "kafka"
	Postgres = "postgres"
	Mysql    = "mysql"
	MongoDB  = "mongo"

	GRPC     = "GRPC"
	SIZE     = "SIZE"
	URI      = "URI"
	STATUS   = "STATUS"
	HTTP     = "HTTP"
	ERROR    = "ERROR"
	METHOD   = "METHOD"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"

	Topic     = "topic"
	Partition = "partition"
	Message   = "message"
	WorkerID  = "workerID"
	Offset    = "offset"
	Time      = "time"

	Page   = "page"
	Size   = "size"
	Search = "search"
	ID     = "id"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

const (
	PermissionReadProfile       = "read_profile"
	PermissionWriteProfile      = "write_profile"
	PermissionManageUsers       = "manage_users"
	PermissionManageRoles       = "manage_roles"
	PermissionManagePermissions = "manage_permissions"
)

var FallbackFutureTime = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
