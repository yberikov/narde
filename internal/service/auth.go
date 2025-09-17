package service

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"narde/internal/transport/http/handlers/model"
)

type AuthService struct {
	tokenGenerator TokenGenerator
}

func NewAuthService(
	tokenGenerator TokenGenerator,
) *AuthService {
	return &AuthService{
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthService) Login(ctx context.Context) (*model.TokenResponse, error) {
	logger := zerolog.Ctx(ctx)

	userID := uuid.Nil

	access, refresh, err := s.tokenGenerator.Tokens(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get tokens")
		return nil, err
	}
	return &model.TokenResponse{
		Access:  access,
		Refresh: refresh,
	}, nil
}
