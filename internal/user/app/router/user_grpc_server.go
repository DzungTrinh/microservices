package router

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/user/domain"
	"microservices/user-management/internal/user/usecases/users"
	"microservices/user-management/pkg/logger"
	userv1 "microservices/user-management/proto/gen/user/v1"
)

type UserGrpcServer struct {
	userv1.UnimplementedUserServiceServer
	uc users.UserUseCase
}

func NewUserGrpcServer(uc users.UserUseCase) *UserGrpcServer {
	return &UserGrpcServer{uc: uc}
}

func (h *UserGrpcServer) Register(ctx context.Context, r *userv1.RegisterRequest) (*userv1.AuthTokens, error) {
	req := domain.RegisterUserReq{
		Email:    r.GetEmail(),
		Username: r.GetUsername(),
		Password: r.GetPassword(),
	}
	resp, err := h.uc.Register(ctx, req)
	if err != nil {
		logger.GetInstance().Errorf("Register failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &userv1.AuthTokens{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		MfaRequired:  resp.MfaRequired,
	}, nil
}

func (h *UserGrpcServer) Login(ctx context.Context, r *userv1.LoginRequest) (*userv1.AuthTokens, error) {
	req := domain.LoginReq{
		Email:    r.GetEmail(),
		Password: r.GetPassword(),
	}
	resp, err := h.uc.Login(ctx, req)
	if err != nil {
		logger.GetInstance().Errorf("Login failed: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &userv1.AuthTokens{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		MfaRequired:  resp.MfaRequired,
	}, nil
}

func (h *UserGrpcServer) Refresh(ctx context.Context, r *userv1.RefreshRequest) (*userv1.AccessToken, error) {
	resp, err := h.uc.RefreshToken(ctx, r.RefreshToken)
	if err != nil {
		logger.GetInstance().Errorf("Refresh failed: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &userv1.AccessToken{AccessToken: resp.AccessToken}, nil
}

func (h *UserGrpcServer) GetAllUsers(ctx context.Context, _ *userv1.Empty) (*userv1.UserList, error) {
	users, err := h.uc.GetAllUsers(ctx)
	if err != nil {
		logger.GetInstance().Errorf("GetAllUsers failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	userList := &userv1.UserList{Users: make([]*userv1.User, len(users))}
	for i, user := range users {
		roles := make([]string, len(user.Roles))
		for j, role := range user.Roles {
			roles[j] = string(role)
		}
		userList.Users[i] = &userv1.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Roles:    roles,
		}
	}
	return userList, nil
}

func (h *UserGrpcServer) GetUserByID(ctx context.Context, r *userv1.GetUserByIDRequest) (*userv1.User, error) {
	user, err := h.uc.GetUserByID(ctx, r.Id)
	if err != nil {
		logger.GetInstance().Errorf("GetUserByID failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = string(role)
	}
	return &userv1.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}, nil
}

func (h *UserGrpcServer) GetCurrentUser(ctx context.Context, _ *userv1.Empty) (*userv1.User, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		logger.GetInstance().Errorf("Missing user_id in context")
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id")
	}

	user, err := h.uc.GetCurrentUser(ctx, userID)
	if err != nil {
		logger.GetInstance().Errorf("GetCurrentUser failed: %v", err)
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = string(role)
	}
	return &userv1.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}, nil
}

func (h *UserGrpcServer) UpdateUserRoles(ctx context.Context, r *userv1.UpdateUserRolesRequest) (*userv1.User, error) {
	user, err := h.uc.UpdateUserRoles(ctx, r.Id, r.Roles)
	if err != nil {
		logger.GetInstance().Errorf("UpdateUserRoles failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = string(role)
	}
	return &userv1.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}, nil
}
