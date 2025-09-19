package service

import (
	"context"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"narde/internal/domain"
	"narde/internal/transport/http/handlers/model"
)

type AuthService struct {
	authRepository AuthRepository
	tokenGenerator TokenGenerator
}

func NewAuthService(
	authRepository AuthRepository,
	tokenGenerator TokenGenerator,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) error {
	logger := zerolog.Ctx(ctx)

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Send()
		return err
	}

	user.Password = string(password)

	if err := s.authRepository.CreateUser(ctx, user); err != nil {
		logger.Error().Err(err).Send()
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.TokenResponse, error) {
	logger := zerolog.Ctx(ctx)

	logger.Info().
		Str("username", req.Username).
		Str("password_from_request", req.Password).
		Msg("Attempting password comparison")

	user, err := s.authRepository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, domain.ErrInvalidCredentials
	}

	logger.Info().
		Str("username", req.Username).
		Str("password_from_request", req.Password).
		Str("hashed_password_from_db", user.Password).
		Msg("Attempting password comparison")

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, domain.ErrInvalidCredentials
	}

	access, refresh, err := s.tokenGenerator.Tokens(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get tokens")
		return nil, err
	}
	return &model.TokenResponse{
		Access:  access,
		Refresh: refresh,
	}, nil
}
