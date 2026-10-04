package service

import (
	"context"
	"errors"
	"net/mail"

	"github.com/MaksimCpp/auth/internal/domain"
	"github.com/MaksimCpp/auth/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type PostgreSQLUserService struct {
	adminCode    string
	repo         domain.UserRepository
	tokenService domain.TokenService
}

func NewPostgreSQLUserService(
	adminCode string,
	repo domain.UserRepository,
	tokenService domain.TokenService,
) *PostgreSQLUserService {
	return &PostgreSQLUserService{
		adminCode: adminCode,
		repo: repo,
		tokenService: tokenService,
	}
}

func (s *PostgreSQLUserService) RegisterUser(
	ctx context.Context, in *dto.RegisterInput,
) (*dto.RegisterOutput, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)

	if err != nil {
		return nil, err
	}

	_, err = mail.ParseAddress(in.Email)

	if err != nil {
		return nil, domain.ErrInvalidEmail
	}

	user := domain.User{
		Username: in.Username,
		Email: in.Email,
		Role: domain.RoleUser,
		PasswordHash: string(passwordHash),
	}

	result, err := s.repo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	out := dto.RegisterOutput{
		UserID: result.ID,
		Email: result.Email,
	}

	return &out, nil
}

func (s *PostgreSQLUserService) RegisterAdmin(
	ctx context.Context, in *dto.RegisterAdminInput,
) (*dto.RegisterOutput, error) {
	if in.AdminCode != s.adminCode {
		return nil, domain.ErrInvalidAdminCode
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)

	if err != nil {
		return nil, err
	}

	_, err = mail.ParseAddress(in.Email)

	if err != nil {
		return nil, domain.ErrInvalidEmail
	}

	user := domain.User{
		Username: in.Username,
		Email: in.Email,
		Role: domain.RoleAdmin,
		PasswordHash: string(passwordHash),
	}

	result, err := s.repo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	out := dto.RegisterOutput{
		UserID: result.ID,
		Email: result.Email,
	}

	return &out, nil
}

func (s *PostgreSQLUserService) RegisterCourier(
	ctx context.Context, in *dto.RegisterInput,
) (*dto.RegisterOutput, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)

	if err != nil {
		return nil, err
	}

	_, err = mail.ParseAddress(in.Email)

	if err != nil {
		return nil, domain.ErrInvalidEmail
	}

	user := domain.User{
		Username: in.Username,
		Email: in.Email,
		Role: domain.RoleCourier,
		PasswordHash: string(passwordHash),
	}

	result, err := s.repo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	out := dto.RegisterOutput{
		UserID: result.ID,
		Email: result.Email,
	}

	return &out, nil
}

func (s *PostgreSQLUserService) LoginUser(
	ctx context.Context, in *dto.LoginInput,
) (*dto.LoginOutput, error) {
	user, err := s.repo.GetByEmail(ctx, in.Email)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))

	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	out := dto.LoginOutput{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &out, nil
}

func (s *PostgreSQLUserService) LoginAdmin(
	ctx context.Context, in *dto.LoginAdminInput,
) (*dto.LoginOutput, error) {
	if in.AdminCode != s.adminCode {
		return nil, domain.ErrInvalidAdminCode
	}
	user, err := s.repo.GetByEmail(ctx, in.Email)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))

	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	out := dto.LoginOutput{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &out, nil
}

func (s *PostgreSQLUserService) LoginCourier(
	ctx context.Context, in *dto.LoginInput,
) (*dto.LoginOutput, error) {
	user, err := s.repo.GetByEmail(ctx, in.Email)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))

	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Role)

	if err != nil {
		return nil, err
	}

	out := dto.LoginOutput{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &out, nil
}

func (s *PostgreSQLUserService) Refresh(
	ctx context.Context, in *dto.RefreshInput,
) (*dto.RefreshOutput, error) {
	userID, role, err := s.tokenService.ValidateRefreshToken(in.RefreshToken)

	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(userID, role)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.tokenService.GenerateAccessToken(userID, role)
	if err != nil {
		return nil, err
	}

	out := dto.RefreshOutput{
		AccessToken: newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return &out, nil
}