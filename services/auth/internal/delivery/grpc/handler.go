package deliverygrpc

import (
	"context"
	"errors"

	"github.com/MaksimCpp/auth/internal/domain"
	"github.com/MaksimCpp/auth/internal/usecase"
	authpb "github.com/MaksimCpp/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer

	registerUC usecase.RegisterUseCase
}

func NewAuthHandler(
	registerUC usecase.RegisterUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
	}
}

func (h *AuthHandler) Register(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {
	in := usecase.RegisterInput{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
	}

	result, err := h.registerUC.Execute(ctx, &in)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExist):
			return nil, status.Error(codes.AlreadyExists, "User already exist.")

		case errors.Is(err, domain.ErrInvalidEmail):
			return nil, status.Error(codes.InvalidArgument, "Invalid email.")
		
		default:
			return nil, status.Error(codes.Internal, "Internal server.")
		}
	}

	response := authpb.RegisterResponse{
		UserId: result.ID,
		Email: result.Email,
	}

	return &response, nil
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	return nil, status.Error(codes.OK, "OK")
}

func (h *AuthHandler) Refresh(
	ctx context.Context,
	req *authpb.RefreshRequest,
) (*authpb.RefreshResponse, error) {
	return nil, status.Error(codes.OK, "OK")
}