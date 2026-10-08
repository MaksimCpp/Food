package deliverygrpc

import (
	"context"
	"errors"

	"github.com/MaksimCpp/auth/internal/domain"
	"github.com/MaksimCpp/auth/internal/dto"
	authpb "github.com/MaksimCpp/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer

	service      domain.UserService
	tokenService domain.TokenService
}

func NewAuthHandler(
	service domain.UserService,
	tokenService domain.TokenService,
) *AuthHandler {
	return &AuthHandler{
		service:      service,
		tokenService: tokenService,
	}
}

func (h *AuthHandler) RegisterUser(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {
	in := dto.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.service.RegisterUser(ctx, &in)

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
		UserId: result.UserID,
		Email:  result.Email,
	}

	return &response, nil
}

func (h *AuthHandler) RegisterAdmin(
	ctx context.Context,
	req *authpb.RegisterAdminRequest,
) (*authpb.RegisterResponse, error) {
	in := dto.RegisterAdminInput{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		AdminCode: req.AdminCode,
	}

	result, err := h.service.RegisterAdmin(ctx, &in)

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
		UserId: result.UserID,
		Email:  result.Email,
	}

	return &response, nil
}

func (h *AuthHandler) RegisterCourier(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {
	in := dto.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.service.RegisterCourier(ctx, &in)

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
		UserId: result.UserID,
		Email:  result.Email,
	}

	return &response, nil
}

func (h *AuthHandler) LoginUser(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	in := dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.service.LoginUser(ctx, &in)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials.")

		default:
			return nil, status.Error(codes.Internal, "Internal server.")
		}
	}

	response := authpb.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	return &response, nil
}

func (h *AuthHandler) LoginAdmin(
	ctx context.Context,
	req *authpb.LoginAdminRequest,
) (*authpb.LoginResponse, error) {
	in := dto.LoginAdminInput{
		Email:     req.Email,
		Password:  req.Password,
		AdminCode: req.AdminCode,
	}

	result, err := h.service.LoginAdmin(ctx, &in)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials.")

		case errors.Is(err, domain.ErrInvalidAdminCode):
			return nil, status.Error(codes.PermissionDenied, "Invalid admin code.")

		default:
			return nil, status.Error(codes.Internal, "Internal server.")
		}
	}

	response := authpb.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	return &response, nil
}

func (h *AuthHandler) LoginCourier(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	in := dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.service.LoginCourier(ctx, &in)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials.")

		default:
			return nil, status.Error(codes.Internal, "Internal server.")
		}
	}

	response := authpb.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	return &response, nil
}

func (h *AuthHandler) Refresh(
	ctx context.Context,
	req *authpb.RefreshRequest,
) (*authpb.RefreshResponse, error) {
	in := dto.RefreshInput{
		RefreshToken: req.RefreshToken,
	}

	result, err := h.service.Refresh(ctx, &in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid refresh token.")
	}

	res := authpb.RefreshResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	return &res, nil
}

func (h *AuthHandler) ValidateAccessToken(
	ctx context.Context,
	req *authpb.ValidateRequest,
) (*authpb.ValidateResponse, error) {
	userID, role, err := h.tokenService.ValidateAccessToken(req.AccessToken)

	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated, "Invalid access token.",
		)
	}

	return &authpb.ValidateResponse{
		UserId: userID,
		Role:   role,
	}, nil
}
