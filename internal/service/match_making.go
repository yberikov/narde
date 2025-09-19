package service

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"narde/internal/domain"
)

type MatchmakingService struct {
	waitingPlayer uuid.UUID
}

func (s *MatchmakingService) FindOpponent(ctx context.Context, playerID uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	logger := zerolog.Ctx(ctx)

	logger.Info().Msgf("Client is looking for a match: %s", playerID)

	if s.waitingPlayer != uuid.Nil && s.waitingPlayer != playerID {
		logger.Info().Msgf("Pairing %s with %s", s.waitingPlayer, playerID)

		s.waitingPlayer = uuid.Nil
		return s.waitingPlayer, playerID, nil
	}

	s.waitingPlayer = playerID
	return uuid.Nil, uuid.Nil, domain.ErrOpponentNotFound
}
