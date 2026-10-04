package domain

import (
	"context"

	"github.com/MaksimCpp/auth/internal/dto"
)

type UserService interface {
	RegisterUser(ctx context.Context, in *dto.RegisterInput) (*dto.RegisterOutput, error)
	RegisterAdmin(ctx context.Context, in *dto.RegisterAdminInput) (*dto.RegisterOutput, error)
	RegisterCourier(ctx context.Context, in *dto.RegisterInput) (*dto.RegisterOutput, error)
	LoginUser(ctx context.Context, in *dto.LoginInput) (*dto.LoginOutput, error)
	LoginAdmin(ctx context.Context, in *dto.LoginAdminInput) (*dto.LoginOutput, error)
	LoginCourier(ctx context.Context, in *dto.LoginInput) (*dto.LoginOutput, error)
	Refresh(ctx context.Context, in *dto.RefreshInput) (*dto.RefreshOutput, error)
}