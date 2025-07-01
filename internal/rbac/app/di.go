package app

import (
	"database/sql"
	"microservices/user-management/cmd/rbac/config"
	"microservices/user-management/internal/rbac/app/router"
	"microservices/user-management/internal/rbac/events"
	"microservices/user-management/internal/rbac/events/handlers"
	"microservices/user-management/internal/rbac/infras/repo"
	"microservices/user-management/internal/rbac/usecases/permission"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/internal/rbac/usecases/role_permission"
	"microservices/user-management/internal/rbac/usecases/user_permission"
	"microservices/user-management/internal/rbac/usecases/user_role"
	"microservices/user-management/pkg/logger"
	"microservices/user-management/pkg/mysql"
)

type Dependencies struct {
	RoleUC          role.RoleUseCase
	PermUC          permission.PermissionUseCase
	UserRoleUC      user_role.UserRoleUseCase
	UserPermUC      user_permission.UserPermissionUseCase
	RolePermUC      role_permission.RolePermissionUseCase
	RBACGrpcHandler *router.RBACGrpcServer
	DB              *sql.DB
}

func InitializeDependencies(cfg config.Config) *Dependencies {
	l := logger.GetInstance()
	l.WithName("rbac-service")

	db, err := mysql.NewDB(cfg.DatabaseDSN)
	if err != nil {
		l.Fatalf("Failed to connect to database: %v", err)
	}

	// Repositories
	roleRepo := repo.NewRoleRepository(db)
	permRepo := repo.NewPermissionRepository(db)
	userRoleRepo := repo.NewUserRoleRepository(db)
	userPermRepo := repo.NewUserPermissionRepository(db)
	rolePermRepo := repo.NewRolePermissionRepository(db)

	// Use cases
	roleUC := role.NewRoleService(roleRepo)
	permUC := permission.NewPermissionService(permRepo)
	userRoleUC := user_role.NewUserRoleService(userRoleRepo, roleUC)
	userPermUC := user_permission.NewUserPermissionService(userPermRepo)
	rolePermUC := role_permission.NewRolePermissionService(rolePermRepo)

	// Consumer
	consumer, err := events.NewConsumer()
	if err != nil {
		logger.GetInstance().Errorf("Failed to initialize RabbitMQ consumer: %v", err)
		panic(err)
	}

	// Register event handlers
	consumer.RegisterHandler("UserRegistered", handlers.NewUserRegisteredHandler(userRoleUC, roleUC))
	consumer.RegisterHandler("AdminUserCreated", handlers.NewAdminUserCreatedHandler(userRoleUC, roleUC))

	rbacGrpcHandler := router.NewRBACGrpcServer(roleUC, permUC, userRoleUC, userPermUC, rolePermUC)

	return &Dependencies{
		RoleUC:          roleUC,
		PermUC:          permUC,
		UserRoleUC:      userRoleUC,
		UserPermUC:      userPermUC,
		RolePermUC:      rolePermUC,
		RBACGrpcHandler: rbacGrpcHandler,
		DB:              db,
	}
}
