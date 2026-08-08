package deliverygrpc

import (
	"context"

	authpb "github.com/MaksimCpp/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(
	ctx context.Context,
	in *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {
	return nil, status.Error(codes.OK, "OK")
}

func (h *AuthHandler) Login(
	ctx context.Context,
	in *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	return nil, status.Error(codes.OK, "OK")
}

func (h *AuthHandler) Refresh(
	ctx context.Context,
	in *authpb.RefreshRequest,
) (*authpb.RefreshResponse, error) {
	return nil, status.Error(codes.OK, "OK")
}