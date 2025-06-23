package router

import (
	"context"
	"microservices/user-management/internal/user/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "microservices/user-management/proto/gen"

	"microservices/user-management/internal/user/usecases/users"
)

type UserGrpcServer struct {
	pb.UnimplementedAuthServiceServer
	uc users.UserUseCase
}

func NewUserGrpcServer(uc users.UserUseCase) *UserGrpcServer {
	return &UserGrpcServer{uc: uc}
}

func (h *UserGrpcServer) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.AuthTokens, error) {
	req := domain.RegisterUserReq{
		Email:    r.GetEmail(),
		Username: r.GetUsername(),
		Password: r.GetPassword(),
	}
	resp, err := h.uc.Register(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AuthTokens{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		MfaRequired:  resp.MfaRequired,
	}, nil
}

func (h *UserGrpcServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.AuthTokens, error) {
	req := domain.LoginReq{
		Email:    r.GetEmail(),
		Password: r.GetPassword(),
	}
	resp, err := h.uc.Login(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &pb.AuthTokens{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		MfaRequired:  resp.MfaRequired,
	}, nil
}

func (h *UserGrpcServer) Refresh(ctx context.Context, r *pb.RefreshRequest) (*pb.AccessToken, error) {
	access, err := h.uc.RefreshToken(ctx, r.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &pb.AccessToken{AccessToken: access.AccessToken}, nil
}
