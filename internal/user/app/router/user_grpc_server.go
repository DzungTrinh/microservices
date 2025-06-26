package router

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices/user-management/internal/user/usecases/auth"
	"microservices/user-management/internal/user/usecases/refresh_token"
	"microservices/user-management/pkg/logger"
	userv1 "microservices/user-management/proto/gen/user/v1"
)

type UserGrpcServer struct {
	userv1.UnimplementedUserServiceServer
	authUC         auth.AuthUseCase
	refreshTokenUC refresh_token.RefreshTokenUseCase
}

func NewUserGrpcServer(authUC auth.AuthUseCase, refreshTokenUC refresh_token.RefreshTokenUseCase) *UserGrpcServer {
	return &UserGrpcServer{
		authUC:         authUC,
		refreshTokenUC: refreshTokenUC,
	}
}

func (h *UserGrpcServer) Register(ctx context.Context, r *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	user, accessToken, refreshToken, err := h.authUC.Register(ctx, r.GetEmail(), r.GetUsername(), r.GetPassword(), r.GetUserAgent(), r.GetIpAddress())
	if err != nil {
		logger.GetInstance().Errorf("Register failed: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &userv1.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		MfaRequired:  false,
	}, nil
}

func (h *UserGrpcServer) Login(ctx context.Context, r *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, accessToken, refreshToken, err := h.authUC.Login(ctx, r.GetEmail(), r.GetPassword(), r.GetUserAgent(), r.GetIpAddress())
	if err != nil {
		logger.GetInstance().Errorf("Login failed: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &userv1.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		MfaRequired:  false,
	}, nil
}

func (h *UserGrpcServer) Refresh(ctx context.Context, r *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := h.refreshTokenUC.RefreshToken(ctx, r.GetRefreshToken(), r.GetUserAgent(), r.GetIpAddress())
	if err != nil {
		logger.GetInstance().Errorf("Refresh failed: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &userv1.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		MfaRequired:  false,
	}, nil
}

//func (h *UserGrpcServer) GetAllUsers(ctx context.Context, _ *userv1.Empty) (*userv1.UserList, error) {
//	users, err := h.uc.GetAllUsers(ctx)
//	if err != nil {
//		logger.GetInstance().Errorf("GetAllUsers failed: %v", err)
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//
//	userList := &userv1.UserList{Users: make([]*userv1.User, len(users))}
//	for i, user := range users {
//		roles := make([]string, len(user.Roles))
//		for j, role := range user.Roles {
//			roles[j] = string(role)
//		}
//		userList.Users[i] = &userv1.User{
//			Id:       user.ID,
//			Username: user.Username,
//			Email:    user.Email,
//			Roles:    roles,
//		}
//	}
//	return userList, nil
//}
//
//func (h *UserGrpcServer) GetUserByID(ctx context.Context, r *userv1.GetUserByIDRequest) (*userv1.User, error) {
//	user, err := h.uc.GetUserByID(ctx, r.Id)
//	if err != nil {
//		logger.GetInstance().Errorf("GetUserByID failed: %v", err)
//		return nil, status.Errorf(codes.NotFound, err.Error())
//	}
//
//	roles := make([]string, len(user.Roles))
//	for i, role := range user.Roles {
//		roles[i] = string(role)
//	}
//	return &userv1.User{
//		Id:       user.ID,
//		Username: user.Username,
//		Email:    user.Email,
//		Roles:    roles,
//	}, nil
//}
//
//func (h *UserGrpcServer) GetCurrentUser(ctx context.Context, _ *userv1.Empty) (*userv1.User, error) {
//	userID, ok := ctx.Value("user_id").(string)
//	if !ok {
//		logger.GetInstance().Errorf("Missing user_id in context")
//		return nil, status.Errorf(codes.Unauthenticated, "missing user_id")
//	}
//
//	user, err := h.uc.GetCurrentUser(ctx, userID)
//	if err != nil {
//		logger.GetInstance().Errorf("GetCurrentUser failed: %v", err)
//		return nil, status.Errorf(codes.NotFound, err.Error())
//	}
//
//	roles := make([]string, len(user.Roles))
//	for i, role := range user.Roles {
//		roles[i] = string(role)
//	}
//	return &userv1.User{
//		Id:       user.ID,
//		Username: user.Username,
//		Email:    user.Email,
//		Roles:    roles,
//	}, nil
//}
//
//func (h *UserGrpcServer) UpdateUserRoles(ctx context.Context, r *userv1.UpdateUserRolesRequest) (*userv1.User, error) {
//	user, err := h.uc.UpdateUserRoles(ctx, r.Id, r.Roles)
//	if err != nil {
//		logger.GetInstance().Errorf("UpdateUserRoles failed: %v", err)
//		return nil, status.Errorf(codes.Internal, err.Error())
//	}
//
//	roles := make([]string, len(user.Roles))
//	for i, role := range user.Roles {
//		roles[i] = string(role)
//	}
//	return &userv1.User{
//		Id:       user.ID,
//		Username: user.Username,
//		Email:    user.Email,
//		Roles:    roles,
//	}, nil
//}
