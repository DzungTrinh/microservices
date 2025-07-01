package handlers

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"microservices/user-management/internal/rbac/domain"
	"microservices/user-management/internal/rbac/events"
	"microservices/user-management/internal/rbac/usecases/role"
	"microservices/user-management/internal/rbac/usecases/user_role"
	"microservices/user-management/pkg/constants"
	"microservices/user-management/pkg/logger"
	"time"
)

type AdminUserCreatedHandler struct {
	urUC       user_role.UserRoleUseCase
	rUC        role.RoleUseCase
	logger     *logger.LoggerService
	maxRetries int
	retryDelay time.Duration
}

func NewAdminUserCreatedHandler(urUC user_role.UserRoleUseCase, rUC role.RoleUseCase) *AdminUserCreatedHandler {
	return &AdminUserCreatedHandler{
		urUC:       urUC,
		rUC:        rUC,
		logger:     logger.GetInstance(),
		maxRetries: 5,
		retryDelay: 1 * time.Second,
	}
}

func (h *AdminUserCreatedHandler) Handle(ctx context.Context, msg amqp091.Delivery) error {
	var payload events.AdminUserCreatedPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		h.logger.Errorf("Failed to unmarshal AdminUserCreated event: %v", err)
		return err
	}

	role, err := h.rUC.GetRoleByName(ctx, constants.RoleAdmin)
	if err != nil {
		h.logger.Errorf("Failed to get role %s: %v", constants.RoleAdmin, err)
		return err
	}

	for attempt := 1; attempt <= h.maxRetries; attempt++ {
		err = h.urUC.AssignRolesToUser(ctx, []domain.UserRole{
			{
				UserID: payload.UserID,
				RoleID: role.ID,
			},
		})
		if err == nil {
			h.logger.Infof("Assigned role %s to user %s", constants.RoleAdmin, payload.UserID)
			return nil
		}
		h.logger.Errorf("Attempt %d/%d: Failed to assign role %s to user %s: %v", attempt, h.maxRetries, constants.RoleAdmin, payload.UserID, err)
		time.Sleep(h.retryDelay)
	}
	return err
}
